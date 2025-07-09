package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ncruces/zenity"
	"log"
)

func openFileDialog() (string, error) {
	filePath, err := zenity.SelectFile()
	if err != nil {
		return "", err
	}
	return filePath, nil
}

type FileForm struct {
	svc            *cliservice.Service
	metadata       textinput.Model
	filePath       string
	isFileSelected bool
}

func NewFileFormUpload(svc *cliservice.Service) FileForm {
	metadata := textinput.New()
	metadata.Placeholder = "Metadata"

	return FileForm{
		svc:            svc,
		metadata:       metadata,
		filePath:       "",
		isFileSelected: false,
	}
}

func (m FileForm) Init() tea.Cmd {
	return nil
}

func (m FileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			if m.isFileSelected {

				fileData := models.FileData{
					Path:     m.filePath,
					Metadata: m.metadata.Value(),
				}
				err := m.svc.AddItem(fileData)
				if err != nil {
					log.Println("Error adding file:", err)
				}
				return NewMainModel(m.svc), nil
			} else {
				m.filePath, err = openFileDialog()
				if err == nil {
					m.isFileSelected = true
				}

				return m, nil
			}
		case "tab":
			m.metadata.Focus()
		}
	}

	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m FileForm) View() string {

	if m.isFileSelected {
		return fmt.Sprintf(
			"File Path: %s\nMetadata: %s\nPress [Enter] to submit\n",
			m.filePath,
			m.metadata.View(),
		)
	} else {
		return "Please select a file first.\n"
	}
}
