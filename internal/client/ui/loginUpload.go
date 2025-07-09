package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type LoginForm struct {
	svc         *cliservice.Service
	serviceName textinput.Model
	login       textinput.Model
	password    textinput.Model
	metadata    textinput.Model
}

func NewLoginFormUpload(svc *cliservice.Service) LoginForm {
	serviceName := textinput.New()
	serviceName.Focused()
	login := textinput.New()
	password := textinput.New()
	metadata := textinput.New()

	serviceName.Placeholder = "Service Name"
	login.Placeholder = "Login"
	password.Placeholder = "Password"
	metadata.Placeholder = "Metadata"

	return LoginForm{
		svc:         svc,
		serviceName: serviceName,
		login:       login,
		password:    password,
		metadata:    metadata,
	}
}

func (m LoginForm) Init() tea.Cmd {
	return nil
}

func (m LoginForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":

			loginData := models.LoginData{
				ServiceName:  m.serviceName.Value(),
				Login:        m.login.Value(),
				PasswordHash: m.password.Value(),
				Metadata:     m.metadata.Value(),
			}
			err := m.svc.AddItem(loginData)
			if err != nil {
				log.Println("Error adding login:", err)
			}
			return NewMainModel(m.svc), nil
		case "tab":

			if m.serviceName.Focused() {
				m.serviceName.Blur()
				m.login.Focus()
			} else if m.login.Focused() {
				m.login.Blur()
				m.password.Focus()
			} else if m.password.Focused() {
				m.password.Blur()
				m.metadata.Focus()
			} else {
				m.metadata.Blur()
				m.serviceName.Focus()
			}
		}
	}

	m.serviceName, _ = m.serviceName.Update(msg)
	m.login, _ = m.login.Update(msg)
	m.password, _ = m.password.Update(msg)
	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m LoginForm) View() string {
	return fmt.Sprintf(
		"ServiceName: %s\nLogin: %s\nPassword: %s\nMetadata: %s\nPress [Enter] to submit\n",
		m.serviceName.View(),
		m.login.View(),
		m.password.View(),
		m.metadata.View(),
	)
}
