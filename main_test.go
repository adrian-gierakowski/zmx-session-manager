package main

import (
	"testing"

	"github.com/mdsakalu/zmx-session-manager/internal/tui"
)

func TestRunSessionManagerReturnsToUIAfterDetach(t *testing.T) {
	requests := []tui.AttachRequest{
		{Target: "demo", Mode: tui.AttachAndReturn},
		{},
	}
	uiCalls := 0
	runUI := func() (tui.AttachRequest, error) {
		request := requests[uiCalls]
		uiCalls++
		return request, nil
	}

	var attached []tui.AttachRequest
	attach := func(request tui.AttachRequest) error {
		attached = append(attached, request)
		return nil
	}

	if err := runSessionManager(runUI, attach); err != nil {
		t.Fatalf("runSessionManager() error = %v", err)
	}
	if uiCalls != 2 {
		t.Fatalf("TUI ran %d times, want 2", uiCalls)
	}
	if len(attached) != 1 || attached[0] != requests[0] {
		t.Fatalf("attach requests = %+v, want %+v", attached, requests[:1])
	}
}
