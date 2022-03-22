# GetPodLogsEfficiently
Purpose of this repo is to demonstrate how to use go routines to get logs from several pod.

In my use case I want to preform an action and then see its output in a pod logs.
I then complicate it with two demands:
1. I want to preform the action multiple times, and check the logs multiple times.
2. I want to have many pods that will react the same way.

Simple get Logs:
[Simple getlogs](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/basicGetLogs/basic.go)

In order to meet the first demand we can poll the logs stream as in the stream example:
[stream example](https://github.com/nwaizer/GetPodLogsEfficiently/blob/main/streamGetLogs/stream.go)