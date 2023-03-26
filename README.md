# GetPodLogsEfficiently
Imagine you have a kubernetes cluster with several pods printing logs.
This is the tale of how I use the tools provided in Go, to fetch the logs from the pods, in a different way.

The original use-case, had many pods listening to a common input and print out a common log line.

I then complicate it with two demands:
1. The common input may change multiple times, and then the code need to check the log multiple times.
2. I want to have many pods that will react the same way.

Simple get Logs:
[Simple getlogs](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/cmd/basicgetlogs/basic.go)

In order to meet the first demand we can poll the logs stream as in the stream example:
[stream example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/cmd/streamgetlogs/stream.go)


To meet the second demand in a more efficient way I will use go-routines:
[go routine example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/cmd/goroutine/streamWgo.go)


And then I add verification and a timeout to the go routines:
[Channeled go routines example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/cmd/channeledgoroutine/channeledgoroutine.go)

# Setup
1. Have the kubeconfig for the cluster in your ~/.kube/config or have a shell environment KUBECONFIG point to it: `$ export KUBECONFIG="~/file"`
2. Deploy the demo pod: `$kubectl apply -f manifests/demo.yaml`
3. Now you can run each of the examples, and get different results.
For example 4, the pod prints 10 good messages and 10 bed, every run will give different results.