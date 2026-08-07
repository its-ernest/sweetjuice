package state

import (
	"github.com/sweet-juice/sweetjuice/app"
)

// MainAppState holds the global application state.
type MainAppState struct {
	Waiting    bool
	WaitingMsg string
}

// NewMainAppState creates a new MainAppState.
func NewMainAppState() *MainAppState {
	return &MainAppState{}
}

// SetWaiting shows a waiting message.
func (s *MainAppState) SetWaiting(msg string) {
	s.Waiting = true
	s.WaitingMsg = msg
	app.ReRender()
}

// ClearWaiting hides the waiting message.
func (s *MainAppState) ClearWaiting() {
	s.Waiting = false
	s.WaitingMsg = ""
	app.ReRender()
}
