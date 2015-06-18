package controller

import (
	"testing"
)

func TestReceiveSMSReply(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// noop
	t.Logf("Not yet implemented")
}

func TestReceiveCallback(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// noop
	t.Logf("Not yet implemented")
}

func TestReceiveDigits(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// noop
	t.Logf("Not yet implemented")
}

func TestGenerateTwiML(t *testing.T) {
	cl, err := NewTestCL()
	if err != nil {
		t.Fatalf("Failed to create test client: %s", err)
	}
	defer cl.Close()

	// noop
	t.Logf("Not yet implemented")
}
