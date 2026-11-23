//go:build !promptcapture

package main

import appruntime "howl-chat/internal/runtime"

// Production builds install nothing: the recorder does not exist in this binary.
func wrapInferenceCapture(client appruntime.Client) appruntime.Client { return client }
