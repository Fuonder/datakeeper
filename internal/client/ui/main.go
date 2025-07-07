package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type MainModel struct {
	svc           *cliservice.Service
	selectedIndex int // Добавляем индекс выбранного объекта
}

func NewMainModel(svc *cliservice.Service) MainModel {
	return MainModel{
		svc:           svc,
		selectedIndex: -1, // Изначально не выбран объект
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
			// Переход к загрузке нового объекта
			return NewUploadModel(m.svc), nil
		case "e":
			// Переход к редактированию выбранного объекта
			logger.Log.Debug("indx", zap.Any("---", m.selectedIndex))
			if m.selectedIndex >= 0 {
				// Передаем индекс выбранного объекта в EditModel
				return NewEditModel(m.svc, m.selectedIndex), nil
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

	// Логины
	result += "Logins:\n"
	for i, login := range dataList.LoginObjects {
		result += fmt.Sprintf("[%d] ID: %d | Service: %s | Login: %s\n", i, login.ID, login.ServiceName, login.Login)
	}

	// Тексты
	result += "Texts:\n"
	for i, text := range dataList.TextObjects {
		result += fmt.Sprintf("[%d] ID: %d | Text: %.20s...\n", i, text.ID, text.Data)
	}

	// Кредитки
	result += "Credit Cards:\n"
	for i, card := range dataList.CreditCardObjects {
		result += fmt.Sprintf("[%d] ID: %d | Card: %s | Owner: %s\n", i, card.ID, card.CardID, card.OwnerName)
	}

	// Файлы
	result += "Files:\n"
	for i, file := range dataList.FileObjects {
		result += fmt.Sprintf("[%d] ID: %d | File: %s | Type: %s\n", i, file.ID, file.Path, file.FileType)
	}

	// Инструкции
	result += "\nCommands:\n"
	result += "[n] - Новый объект\n"
	result += "[e] - Редактировать объект\n"
	result += "[q] - Выход\n"

	return result
}
