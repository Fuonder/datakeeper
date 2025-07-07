package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type authMode int

const (
	modeLogin authMode = iota
	modeRegister
)

type AuthModel struct {
	login    textinput.Model
	password textinput.Model
	svc      *cliservice.Service
	user     models.User
	authMode authMode // Добавляем поле для текущего режима
	errorMsg string   // Добавляем поле для ошибок
}

func NewAuthModel(svc *cliservice.Service) AuthModel {
	login := textinput.New()
	login.Placeholder = "Login"
	login.Focus()
	login.CharLimit = 64
	login.Width = 20

	password := textinput.New()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword
	password.CharLimit = 64
	password.Width = 20

	return AuthModel{
		login:    login,
		password: password,
		svc:      svc,
		authMode: modeLogin,
	}
}

func (m AuthModel) Init() tea.Cmd {
	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// Выход
			return m, tea.Quit
		case "enter":
			// При нажатии Enter выполняем логин или регистрацию в зависимости от выбранного режима
			m.user.Login = m.login.Value()
			m.user.PwdHash = m.password.Value()

			var err error
			if m.authMode == modeRegister {
				// Регистрация
				err = m.svc.Register(m.user)
			} else if m.authMode == modeLogin {
				// Логин
				err = m.svc.Login(m.user)
			}

			// Если произошла ошибка, возвращаем текущую модель
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}

			// Переход на главный экран после успешного логина или регистрации
			return NewMainModel(m.svc), nil
		case "tab":
			// Переключение между полями login и password
			if m.login.Focused() {
				m.login.Blur()
				m.password.Focus()
			} else {
				m.password.Blur()
				m.login.Focus()
			}
		case "l":
			// Переключаемся на логин
			m.authMode = modeLogin
			m.errorMsg = ""
		case "r":
			// Переключаемся на регистрацию
			m.authMode = modeRegister
			m.errorMsg = ""
		}
	}

	// Обновляем оба поля
	m.login, _ = m.login.Update(msg)
	m.password, _ = m.password.Update(msg)
	return m, nil
}

func (m AuthModel) View() string {
	var modeStr string
	if m.authMode == modeRegister {
		modeStr = "Register"
	} else {
		modeStr = "Login"
	}

	// Строим интерфейс в зависимости от текущего режима
	b := fmt.Sprintf("== %s ==\n\n", modeStr)
	b += fmt.Sprintf("%s\n%s\n\n", m.login.View(), m.password.View())
	b += "Press [Tab] to switch input.\n"
	b += "Press [Enter] to submit.\n"
	b += "Press [l] to login, [r] to register.\n"

	// Показываем ошибку, если она есть
	if m.errorMsg != "" {
		b += fmt.Sprintf("\n[ERROR]: %s\n", m.errorMsg)
	}

	return b
}
