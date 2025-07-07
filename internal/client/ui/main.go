package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

type MainModel struct {
	svc           *cliservice.Service
	selectedIndex int
	selectedType  string
}

func NewMainModel(svc *cliservice.Service) MainModel {
	return MainModel{
		svc:           svc,
		selectedIndex: -1,
		selectedType:  "",
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":

			return NewUploadModel(m.svc), nil
		case "e":

			return NewEditModel(m.svc), nil

		}

		if msg.String() >= "0" && msg.String() <= "9" {
			selectedIndex, err := strconv.Atoi(msg.String())
			if err == nil {

				dataList, err := m.svc.FetchData()
				if err != nil {
					return m, nil
				}

				if selectedIndex < len(dataList.LoginObjects) {
					m.selectedType = "login"
					m.selectedIndex = selectedIndex
				} else if selectedIndex-len(dataList.LoginObjects) < len(dataList.TextObjects) {

					m.selectedType = "text"
					m.selectedIndex = selectedIndex - len(dataList.LoginObjects)
				} else if selectedIndex-len(dataList.LoginObjects)-len(dataList.TextObjects) < len(dataList.CreditCardObjects) {

					m.selectedType = "creditcard"
					m.selectedIndex = selectedIndex - len(dataList.LoginObjects) - len(dataList.TextObjects)
				} else if selectedIndex-len(dataList.LoginObjects)-len(dataList.TextObjects)-len(dataList.CreditCardObjects) < len(dataList.FileObjects) {

					m.selectedType = "file"
					m.selectedIndex = selectedIndex - len(dataList.LoginObjects) - len(dataList.TextObjects) - len(dataList.CreditCardObjects)
				}
			}
		}
	}

	return m, nil
}

func (m MainModel) View() string {
	var result string
	dataList, err := m.svc.FetchData()
	if err != nil {
		return err.Error()
	}

	result += "Logins:\n"
	for i, login := range dataList.LoginObjects {
		result += fmt.Sprintf("[%d] ID: %d | Service: %s | Login: %s\n", i, login.ID, login.ServiceName, login.Login)
	}

	result += "Texts:\n"
	for i, text := range dataList.TextObjects {
		result += fmt.Sprintf("[%d] ID: %d | Text: %.20s...\n", i, text.ID, text.Data)
	}

	result += "Credit Cards:\n"
	for i, card := range dataList.CreditCardObjects {
		result += fmt.Sprintf("[%d] ID: %d | Card: %s | Owner: %s\n", i, card.ID, card.CardID, card.OwnerName)
	}

	result += "Files:\n"
	for i, file := range dataList.FileObjects {
		result += fmt.Sprintf("[%d] ID: %d | File: %s | Type: %s\n", i, file.ID, file.Path, file.FileType)
	}

	result += "\nCommands:\n"
	result += "[n] - Новый объект\n"
	result += "[e] - Редактировать объект\n"
	result += "[q] - Выход\n"

	return result
}
