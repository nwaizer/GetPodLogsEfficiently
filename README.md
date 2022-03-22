# GetPodLogsEfficiently
Imaging you have a kubernetes cluster with several pods printing logs.   
Purpose of this repo is to share several examples for fetching logs from pods.

In my use case I want to preform an action and then see it print to the logs in the pods.
I then complicate it with two demands:
1. I want to preform the action multiple times, and check the logs multiple times.
2. I want to have many pods that will react the same way.

Simple get Logs:
[Simple getlogs](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/basicGetLogs/basic.go)

In order to meet the first demand we can poll the logs stream as in the stream example:
[stream example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/streamGetLogs/stream.go)

To meet the second demand in a more efficient way I will use go-routines:
[go routine example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/3_goRoutine/streamWgo.go)

And then I add verification and a timeout to the go routines:
[Channeled go routines example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/4_channeledGoRoutine/main.go)
