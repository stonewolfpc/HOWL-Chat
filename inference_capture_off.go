//go:build !promptcapture

package main

// Production builds install nothing: the recorder does not exist in this binary.
func wrapInferenceCapture(client inferenceClient) inferenceClient { return client }
