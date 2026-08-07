package state

import (
	"github.com/sweet-juice/sweetjuice/app"
)

type MainAppState struct {
	SelectedTab string
}

func NewMainAppState() *MainAppState {
	return &MainAppState{
		SelectedTab: "for-you",
	}
}

func (s *MainAppState) SetSelectedTab(tab string) {
	s.SelectedTab = tab
	app.ReRender()
}
