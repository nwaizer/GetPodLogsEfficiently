
<table>
<tr>
<th>Go code</th>
<th>Comments</th>
</tr>
<tr>
<td>
<pre>

  func main() {
	PodsList := utils.GetPods()
	ctx := context.Background()
	cancelCtx, endCheckers := context.WithCancel(ctx)

	outChannel := make(chan string, len(PodsList.Items))

	for _, Pod := range PodsList.Items {
		go Checker(cancelCtx, Pod.Name, outChannel)
	}

	if !VerifyEvents(PodsList, outChannel, 10) {
		endCheckers()
		log.Errorf("Test failed due to error printed above")
	}
}


</pre>
</td>
<td>
<pre>
Here we first get a list of pods. each pod runs a single container, generating the log line
Now we will use a context to allow use to stop the go routines that run in the background
.
.
.
.
Each go routine will check the log line and send the result to this channel
Spin-up a go routine for each pod.
.
.
.
Now run this collector function that will ignite in case a single of the go routines fail a verify. 
.
.
.
.
</pre>
</td>
</tr>
</table>