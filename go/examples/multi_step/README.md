### Running the example

```bash
gcloud builds submit --config=./cloudbuild.yaml
```

### Note

This trivial example builds 2 binaries that share a Golang package `github.com/golang/glog`.

Because the cloudbuild.yaml shares a volume called `go-modules` across all steps, once the `glog` package is pulled for the first step (`Step #0`), it is available to be used by the second step. In practice, sharing immutable packages this way should significantly improve build times.

```bash
BUILD
Starting Step #0
Step #0: Pulling image: golang:1.12
Step #0: 1.12: Pulling from library/golang
Step #0: Digest: sha256:017b7708cffe1432adfc8cc27119bfe4c601f369f3ff086c6e62b3c6490bf540
Step #0: Status: Downloaded newer image for golang:1.12
Step #0: go: finding github.com/golang/glog latest
Step #0: go: downloading github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Step #0: go: extracting github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Finished Step #0
Starting Step #1
Step #1: Already have image: golang:1.12
Finished Step #1
Starting Step #2
Step #2: Already have image: busybox
Step #2: total 4840
Step #2: -rwxr-xr-x    1 root     root       2474479 Jul 10 23:59 bar
Step #2: -rwxr-xr-x    1 root     root       2474479 Jul 10 23:59 foo
Finished Step #2
PUSH
DONE
```

If the `options` second were commented out and the `busybox` step removed (as it will fail), the output would show the second step repulling (unnecessarily) the `glog` package:

```bash
BUILD
Starting Step #0
Step #0: Pulling image: golang:1.12
Step #0: 1.12: Pulling from library/golang
Step #0: Digest: sha256:017b7708cffe1432adfc8cc27119bfe4c601f369f3ff086c6e62b3c6490bf540
Step #0: Status: Downloaded newer image for golang:1.12
Step #0: go: finding github.com/golang/glog latest
Step #0: go: downloading github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Step #0: go: extracting github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Finished Step #0
Starting Step #1
Step #1: Already have image: golang:1.12
Step #1: go: finding github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Step #1: go: downloading github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Step #1: go: extracting github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
Finished Step #1
PUSH
DONE

```


This example also uses the Golang team's Module Mirror [https://proxy.golang.org](https://proxy.golang.org) for its performance and security benefits.