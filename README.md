# hello-universe

## Build

This repo is go-gettable and works with the standard go tools. It also works with [Bazel](https://www.bazel.build).

Building with Bazel requires a bazel 5.2.0+ installation and running the following command:

```
bazel build //...
```

After the build is complete the `hello-universe` binary lives under the bazel-bin directory: 

```
bazel-bin/hello-universe
```  

## Example Usage

```
$ hello-universe -h
```
```
Usage of hello-universe:
  -api-server string
    	Kubernetes API server (default "http://127.0.0.1:8080")
  -cpu-limit string
    	Max CPU in milicores (default "100m")
  -cpu-request string
    	Min CPU in milicores (default "100m")
  -http string
    	HTTP service address (default "127.0.0.1:80")
  -kubernetes
    	Deploy to Kubernetes.
  -memory-limit string
    	Max memory in MB (default "64M")
  -memory-request string
    	Min memory in MB (default "64M")
  -replicas int
    	Number of replicas (default 1)
```

### Example

```
$ export HELLO_UNIVERSE_TOKEN="7346053dafaf4a24825790f4389704f5"
```

```
$ hello-universe -kubernetes -replicas 3  -cpu-limit 200m -memory-limit 500M
```

```
$ kubectl get replicasets
```

```
NAME             DESIRED   CURRENT   READY     AGE
hello-universe   3         3         0         15s
```

```
$ kubectl get pods
```
```
NAME                   READY     STATUS    RESTARTS   AGE
hello-universe-2genp   0/1       Pending   0          43s
hello-universe-gg1eu   0/1       Pending   0          43s
hello-universe-zgpuy   0/1       Pending   0          43s
```
