package ui

import (
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	tea "github.com/charmbracelet/bubbletea"
)

// UploadModel для загрузки нового объекта
type UploadModel struct {
	svc        *cliservice.Service
	objectType string // Тип объекта (login, text, creditcard, file)

}

func NewUploadModel(svc *cliservice.Service) UploadModel {

	return UploadModel{
		svc:        svc,
		objectType: "",
	}
}

func (m UploadModel) Init() tea.Cmd {
	return nil
}

// Метод Update обрабатывает данные в зависимости от типа объекта
func (m UploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			m.objectType = "login"
			return NewLoginForm(m.svc), nil
		case "t":
			m.objectType = "text"
			return NewTextForm(m.svc), nil
		case "c":
			m.objectType = "creditcard"
			return NewCreditCardForm(m.svc), nil
		case "f":
			m.objectType = "file"
			return NewFileForm(m.svc), nil
		}
	}

	return m, nil
}

func (m UploadModel) View() string {
	var b string
	b = "Press [l] for Login, [t] for Text, [c] for Credit Card, [f] for File\n"
	return b
}
