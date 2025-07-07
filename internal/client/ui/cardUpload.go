package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type CreditCardForm struct {
	svc       *cliservice.Service
	cardID    textinput.Model
	ownerName textinput.Model
	metadata  textinput.Model
}

func NewCreditCardForm(svc *cliservice.Service) CreditCardForm {
	cardID := textinput.New()
	cardID.Focused()
	ownerName := textinput.New()
	metadata := textinput.New()

	cardID.Placeholder = "Card ID"
	ownerName.Placeholder = "Owner Name"
	metadata.Placeholder = "Metadata"

	return CreditCardForm{
		svc:       svc,
		cardID:    cardID,
		ownerName: ownerName,
		metadata:  metadata,
	}
}

func (m CreditCardForm) Init() tea.Cmd {
	return nil
}

func (m CreditCardForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			// Сохранение данных
			creditCardData := models.CreditCardData{
				CardID:    m.cardID.Value(),
				OwnerName: m.ownerName.Value(),
				Metadata:  m.metadata.Value(),
			}
			err := m.svc.AddItem(creditCardData)
			if err != nil {
				log.Println("Error adding credit card:", err)
			}
			return NewMainModel(m.svc), nil
		case "tab":
			// Переключение между полями login и password
			if m.cardID.Focused() {
				m.cardID.Blur()
				m.ownerName.Focus()
			} else if m.ownerName.Focused() {
				m.ownerName.Blur()
				m.metadata.Focus()
			} else {
				m.metadata.Blur()
				m.cardID.Focus()
			}
		}
	}

	// Обновление полей
	m.cardID, _ = m.cardID.Update(msg)
	m.ownerName, _ = m.ownerName.Update(msg)
	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m CreditCardForm) View() string {
	return fmt.Sprintf(
		"Card ID: %s\nOwner Name: %s\nMetadata: %s\nPress [Enter] to submit\n",
		m.cardID.View(),
		m.ownerName.View(),
		m.metadata.View(),
	)
}
