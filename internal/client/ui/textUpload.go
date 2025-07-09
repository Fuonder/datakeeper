package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type TextForm struct {
	svc      *cliservice.Service
	data     textinput.Model
	metadata textinput.Model
}

func NewTextFormUpload(svc *cliservice.Service) TextForm {
	data := textinput.New()
	data.Focused()
	metadata := textinput.New()

	data.Placeholder = "Text Data"
	metadata.Placeholder = "Metadata"

	return TextForm{
		svc:      svc,
		data:     data,
		metadata: metadata,
	}
}

func (m TextForm) Init() tea.Cmd {
	return nil
}

func (m TextForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":

			textData := models.TextData{
				Data:     m.data.Value(),
				Metadata: m.metadata.Value(),
			}
			err := m.svc.AddItem(textData)
			if err != nil {
				log.Println("Error adding text:", err)
			}
			return NewMainModel(m.svc), nil
		case "tab":

			if m.data.Focused() {
				m.data.Blur()
				m.metadata.Focus()
			} else {
				m.metadata.Blur()
				m.data.Focus()
			}
		}
	}

	m.data, _ = m.data.Update(msg)
	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m TextForm) View() string {
	return fmt.Sprintf(
		"Text Data: %s\nMetadata: %s\nPress [Enter] to submit\n",
		m.data.View(),
		m.metadata.View(),
	)
}
