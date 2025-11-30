package tui

import "sshbuddy/pkg/models"

// ConnectMsg is a Bubble Tea message for SSH connection
type ConnectMsg struct {
	Host models.Host
}

// PingResultMsg is a Bubble Tea message for ping results
type PingResultMsg struct {
	Host     models.Host
	Status   bool   // true if reachable
	PingTime string // ping time in ms
}

// FormSubmittedMsg is sent when a form is submitted
type FormSubmittedMsg struct {
	Host models.Host
}

// TermixAuthSuccessMsg is sent when Termix authentication succeeds
type TermixAuthSuccessMsg struct{}
