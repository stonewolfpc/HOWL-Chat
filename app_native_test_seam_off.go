//go:build !native_test_seam

package main

// Production builds compile only this no-op stub. The test seam and its
// environment switch exist exclusively under the native_test_seam build tag.
func startNativeTestSeam(*App) {}

func stopNativeTestSeam() {}
