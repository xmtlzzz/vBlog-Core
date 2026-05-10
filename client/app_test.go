package main

import "testing"

func TestApp_Connect_InvalidAddress(t *testing.T) {
	app := NewApp()
	err := app.Connect("invalid:99999", "test-key")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestApp_Disconnect(t *testing.T) {
	app := NewApp()
	err := app.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect on unconnected app should not error: %v", err)
	}
}

func TestApp_GetLatestStats_NotConnected(t *testing.T) {
	app := NewApp()
	_, err := app.GetLatestStats()
	if err == nil {
		t.Fatal("expected error when not connected")
	}
}
