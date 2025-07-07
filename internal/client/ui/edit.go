package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditModel struct {
	svc      *cliservice.Service
	objectID int
	input    textinput.Model
	filePath string
}

func NewEditModel(svc *cliservice.Service, objectID int) EditModel {
	input := textinput.New()
	input.Placeholder = "Edit data"
	input.Focus()
	input.CharLimit = 100
	input.Width = 20
	return EditModel{
		svc:      svc,
		objectID: objectID,
		input:    input,
	}
}

func (m EditModel) Init() tea.Cmd {
	return nil
}

func (m EditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			//// Реализуем редактирование объекта, переданного по ID
			//if textData, err := m.storage.GetTextItemByIndex(m.objectID); err == nil {
			//	// Редактируем текст
			//	textData.Data = m.input.Value()
			//	m.storage.AddItem(textData)
			//}
		}
	}
	m.input, _ = m.input.Update(msg)
	return m, nil
}

func (m EditModel) View() string {
	return fmt.Sprintf("Edit object ID: %d\n%s", m.objectID, m.input.View())
}
