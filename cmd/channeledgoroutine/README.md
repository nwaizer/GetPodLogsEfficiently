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
Here we first get a list of pods. each pod runs a single container, generating the log line
.
Now we will use a context to allow use to stop the go routines that run in the background
.
.
Each go routine will check the log line and send the result to this channel
.
Spin-up a go routine for each pod.
.
.
Now run this collector function that will ignite in case a single of the go routines fail a verify.
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
