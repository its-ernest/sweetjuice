package state

import (
	"github.com/sweet-juice/sweetjuice/app"
)

// MainAppState is the single source of truth for the entire app.
type MainAppState struct {
	User *UserState
}

func NewMainAppState() *MainAppState {
	return &MainAppState{
		User: &UserState{Name: ""},
	}
}

// UserState handles user-related data and logic
type UserState struct {
	Name string
}

func (u *UserState) SetName(name string) {
	u.Name = name
	app.ReRender()
}
