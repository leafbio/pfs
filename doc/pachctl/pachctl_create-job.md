## ./pachctl create-job

Create a new job. Returns the id of the created job.

### Synopsis


Create a new job from a spec, the spec looks like this
{
  "transform": {
    "cmd": [
      "cmd",
      "args..."
    ],
    "env": {
      "foo": "bar"
    },
    "secrets": [
      {
        "name": "secret_name",
        "mountPath": "/path/in/container"
      }
    ],
    "imagePullSecrets": [
      "my-secret"
    ],
    "acceptReturnCode": [
      "1"
    ]
  },
  "parallelismSpec": {
    "constant": "1"
  },
  "inputs": [
    {
      "commit": {
        "repo": {
          "name": "in_repo"
        },
        "id": "10cf676b626044f9a405235bf7660959"
      },
      "glob": "*",
      "lazy": true
    }
  ]
}

```
./pachctl create-job -f job.json
```

### Options

```
      --delete-job        delete the jobs in this pipeline as well
  -f, --file string       The file containing the job, it can be a url or local file. - reads from stdin. (default "-")
      --password string   Your password for the registry being pushed to.
  -p, --push-images       If true, push local docker images into the cluster registry.
  -r, --registry string   The registry to push images to. (default "docker.io")
  -u, --username string   The username to push images as, defaults to your OS username.
```

### Options inherited from parent commands

```
      --no-metrics   Don't report user metrics for this command
  -v, --verbose      Output verbose logs
```

### SEE ALSO
* [./pachctl](./pachctl.md)	 - 

###### Auto generated by spf13/cobra on 7-Apr-2017
