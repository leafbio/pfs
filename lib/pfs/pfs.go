package pfs

import (
	"io"
	"net/http"
	"path"
)

func url(name string) string {
	return path.Join("pfs", "pfs", name)
}

func Post(name string, r io.Reader) error {
	http := &http.Client{}
	r, err := http.Post(url(name), "", r)
	return err
}

func Put(name string, r io.Reader) error {
	http := &http.Client{}
	r, err := http.Put(url(name), "", r)
	return err
}

func Get(name string) (io.Reader, error) {
	http := &http.Client{}
	r, err := http.Get(url(name))
	return r.Body, err
}

func GetAt(name, commit string), (io.Reader, error) {
	http := &http.Client{}
	r, err := http.Get(url(name))
	return r.Body, err
}

func Delete(name string) error {
	http := &http.Client{}
	req := http.NewRequest("DELETE", url(name), nil)
	_, err := http.Do(req)
	return err
}
