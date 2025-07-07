package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type CreditCardFormRedact struct {
	svc       *cliservice.Service
	cardID    textinput.Model
	ownerName textinput.Model
	metadata  textinput.Model
	card      models.CreditCardData
}

func NewCardFormRedact(svc *cliservice.Service, obj models.CreditCardData) CreditCardFormRedact {
	cardID := textinput.New()
	cardID.Focused()
	ownerName := textinput.New()
	metadata := textinput.New()

	cardID.Placeholder = "Card ID"
	ownerName.Placeholder = "Owner Name"
	metadata.Placeholder = "Metadata"
	cardID.SetValue(obj.CardID)
	ownerName.SetValue(obj.OwnerName)
	metadata.SetValue(obj.Metadata)

	return CreditCardFormRedact{
		svc:       svc,
		cardID:    cardID,
		ownerName: ownerName,
		metadata:  metadata,
		card:      obj,
	}
}

func (m CreditCardFormRedact) Init() tea.Cmd {
	return nil
}

func (m CreditCardFormRedact) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":

			creditCardData := models.CreditCardData{
				ID:        m.card.ID,
				CardID:    m.cardID.Value(),
				OwnerName: m.ownerName.Value(),
				Metadata:  m.metadata.Value(),
			}
			err := m.svc.AddItem(creditCardData)
			if err != nil {
				log.Println("Error updating credit card:", err)
			}
			return NewMainModel(m.svc), nil
		case "tab":

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

	m.cardID, _ = m.cardID.Update(msg)
	m.ownerName, _ = m.ownerName.Update(msg)
	m.metadata, _ = m.metadata.Update(msg)

	return m, nil
}

func (m CreditCardFormRedact) View() string {
	return fmt.Sprintf(
		"Card ID: %s\nOwner Name: %s\nMetadata: %s\nPress [Enter] to submit\n",
		m.cardID.View(),
		m.ownerName.View(),
		m.metadata.View(),
	)
}
