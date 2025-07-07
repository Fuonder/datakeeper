package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type FileFormRedact struct {
	svc      *cliservice.Service
	metadata textinput.Model
	filePath string
	file     models.FileData
}

func NewFileFormRedact(svc *cliservice.Service, file models.FileData) FileFormRedact {
	metadata := textinput.New()
	metadata.Placeholder = "Metadata"
	metadata.SetValue(file.Metadata)

	return FileFormRedact{
		svc:      svc,
		metadata: metadata,
		filePath: file.Path,
		file:     file,
	}
}

func (m FileFormRedact) Init() tea.Cmd {
	return nil
}

func (m FileFormRedact) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":

			fileData := models.FileData{
				ID:       m.file.ID,
				Path:     m.filePath,
				Metadata: m.metadata.Value(),
			}
			err := m.svc.AddItem(fileData)
			if err != nil {
				log.Println("Error updating file:", err)
			}
			return NewMainModel(m.svc), nil
		case "tab":

			m.metadata.Focus()
		}
	}

	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m FileFormRedact) View() string {
	return fmt.Sprintf(
		"File Path: %s\nMetadata: %s\nPress [Enter] to submit\n",
		m.filePath,
		m.metadata.View(),
	)
}
