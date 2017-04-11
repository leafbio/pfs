package obj

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/storage"
)

type microsoftClient struct {
	blobClient storage.BlobStorageClient
	container  string
}

func newMicrosoftClient(container string, accountName string, accountKey string) (*microsoftClient, error) {
	client, err := storage.NewBasicClient(
		accountName,
		accountKey,
	)
	if err != nil {
		return nil, err
	}

	return &microsoftClient{
		blobClient: client.GetBlobService(),
		container:  container,
	}, nil
}

func (c *microsoftClient) Writer(name string) (io.WriteCloser, error) {
	writer, err := newMicrosoftWriter(c, name)
	if err != nil {
		return nil, err
	}
	return newBackoffWriteCloser(c, writer), nil
}

func (c *microsoftClient) Reader(name string, offset uint64, size uint64) (io.ReadCloser, error) {
	byteRange := byteRange(offset, size)
	var reader io.ReadCloser
	var err error
	if byteRange == "" {
		reader, err = c.blobClient.GetBlob(c.container, name)
	} else {
		reader, err = c.blobClient.GetBlobRange(c.container, name, byteRange, nil)
	}

	if err != nil {
		return nil, err
	}
	return newBackoffReadCloser(c, reader), nil
}

func (c *microsoftClient) Delete(name string) error {
	return c.blobClient.DeleteBlob(c.container, name, nil)
}

func (c *microsoftClient) Walk(name string, fn func(name string) error) error {
	blobList, err := c.blobClient.ListBlobs(c.container, storage.ListBlobsParameters{Prefix: name})
	if err != nil {
		return err
	}

	for _, file := range blobList.Blobs {
		if err := fn(file.Name); err != nil {
			return err
		}
	}
	return nil
}

func (c *microsoftClient) Exists(name string) bool {
	exists, _ := c.blobClient.BlobExists(c.container, name)
	return exists
}

func (c *microsoftClient) isRetryable(err error) (ret bool) {
	microsoftErr, ok := err.(storage.AzureStorageServiceError)
	if !ok {
		return false
	}
	return microsoftErr.StatusCode >= 500
}

func (c *microsoftClient) IsNotExist(err error) bool {
	microsoftErr, ok := err.(storage.AzureStorageServiceError)
	if !ok {
		return false
	}
	return microsoftErr.StatusCode == 404
}

func (c *microsoftClient) IsIgnorable(err error) bool {
	return false
}

type microsoftWriter struct {
	container  string
	blob       string
	blobClient storage.BlobStorageClient
}

func newMicrosoftWriter(client *microsoftClient, name string) (*microsoftWriter, error) {
	// create container
	_, err := client.blobClient.CreateContainerIfNotExists(client.container, storage.ContainerAccessTypePrivate)
	if err != nil {
		return nil, err
	}

	// create blob
	err = client.blobClient.CreateBlockBlob(client.container, name)
	if err != nil {
		return nil, err
	}

	return &microsoftWriter{
		container:  client.container,
		blob:       name,
		blobClient: client.blobClient,
	}, nil
}

func (w *microsoftWriter) Write(b []byte) (int, error) {
	blockList, err := w.blobClient.GetBlockList(w.container, w.blob, storage.BlockListTypeAll)
	if err != nil {
		return 0, err
	}

	blocksLen := len(blockList.CommittedBlocks)
	amendList := []storage.Block{}
	for _, v := range blockList.CommittedBlocks {
		amendList = append(amendList, storage.Block{v.Name, storage.BlockStatusCommitted})
	}

	var chunkSize = storage.MaxBlobBlockSize
	inputSourceReader := bytes.NewReader(b)
	chunk := make([]byte, chunkSize)
	for {
		n, err := inputSourceReader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%011d\n", blocksLen)))
		data := chunk[:n]
		err = w.blobClient.PutBlock(w.container, w.blob, blockID, data)
		if err != nil {
			return 0, err
		}
		// add current uncommitted block to temporary block list
		amendList = append(amendList, storage.Block{blockID, storage.BlockStatusUncommitted})
		blocksLen++
	}

	// update block list to blob committed block list.
	err = w.blobClient.PutBlockList(w.container, w.blob, amendList)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *microsoftWriter) Close() error {
	return nil
}
