# GetPodLogsEfficiently
Imaging you have a kubernetes cluster with several pods printing logs.   
Purpose of this repo is to share several examples for fetching logs from pods.

In my use case I want to preform an action and then see it print to the logs in the pods.
I then complicate it with two demands:
1. I want to preform the action multiple times, and check the logs multiple times.
2. I want to have many pods that will react the same way.

Simple get Logs:
[Simple getlogs](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/1_basicGetLogs/basic.go)

In order to meet the first demand we can poll the logs stream as in the stream example:
[stream example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/2_streamGetLogs/stream.go)


To meet the second demand in a more efficient way I will use go-routines:
[go routine example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/3_goRoutine/streamWgo.go)


And then I add verification and a timeout to the go routines:
[Channeled go routines example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/4_channeledGoRoutine/main.go)

# Setup
1. Have the kubeconfig for the cluster in your ~/.kube/config or have a shell environment KUBECONFIG point to it: `$ export KUBECONFIG="~/file"`
2. Deploy the demo pod: `$kubectl apply -f manifests/demo.yaml`
3. Now you can run each of the examples, and get different results.
For example 4, the pod prints 10 good messages and 10 bed, every run will give different results.