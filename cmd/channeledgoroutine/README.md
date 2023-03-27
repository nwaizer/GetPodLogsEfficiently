<head>
  <title>Line by line</title>
  <link rel=stylesheet href="https://gobyexample.com/site.css">
</head>
<div class="example" id="goroutines">
<table>
<tr>
<th>Go code</th>
<th>Comments</th>
</tr>
<tr>
<td class="docs">
<pre>
  func main() {
	Client = client.New("")

	ctx := context.Background()
	cancelCtx, endCheckers := context.WithCancel(ctx)
    //

	pods, err := client.Client.Pods(Namespace).List(context.Background(), v12.ListOptions{
		LabelSelector: LabelSelector})

	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)

		return
	}

	outChannel := make(chan string, len(pods.Items))

	for _, Pod := range pods.Items {
		go Checker(cancelCtx, Pod.Name, outChannel)
	}

	if !VerifyEvents(pods, outChannel, 10) {
		endCheckers()
		logrus.Errorf("Test failed due to error printed above")
	}
}


</pre>
</td>

<td>
<pre>
.
Here we first get the list of pods, that we will use for this example.
Each pod runs a single container, generating the log line.
.
A context is used to stop the go routines we will use.
.
.
Each go routine will check the log line and send the result to this channel
.
.
.
.
.
.

.
.
.
.
Spin-up a go routine for each pod.
.
.
Now run this collector function that will check all go routines checked the log as expected.
.
</pre>
</td>
.
</table>
</div>
.
.
.
.
