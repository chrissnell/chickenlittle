package controller

import (
	"testing"
)

func TestNotifyPersonViaVictorops(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// noop
	t.Logf("Not yet implemented")
}
