package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/impossibleclone/vimingo/internal/app" // Imports the new app package
)

func main() {
	// Initialize the model from the internal/app package
	m := app.InitialModel()

	// WithAltScreen to have full screen
	p := tea.NewProgram(&m, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
