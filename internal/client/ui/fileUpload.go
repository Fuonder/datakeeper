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

// Функция для открытия диалога выбора файла через zenity
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
	isFileSelected bool // Флаг, чтобы знать, что файл уже выбран
}

func NewFileForm(svc *cliservice.Service) FileForm {
	metadata := textinput.New()
	metadata.Placeholder = "Metadata"

	return FileForm{
		svc:            svc,
		metadata:       metadata,
		filePath:       "",
		isFileSelected: false, // Изначально файл не выбран
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
				// Сохранение данных, если файл был выбран
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
				// Если файл не выбран, просим выбрать файл
				return m, nil
			}
		case "tab":
			m.metadata.Focus()
		}
	}

	// Обновление полей
	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m FileForm) View() string {
	// Отображаем путь к файлу, если он выбран, или сообщение о необходимости выбрать файл
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
