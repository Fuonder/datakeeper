package ui

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
)

type EditModel struct {
	svc        *cliservice.Service
	objectType string
}

func NewEditModel(svc *cliservice.Service) EditModel {
	return EditModel{
		svc:        svc,
		objectType: "",
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
		case "l":
			m.objectType = "login"
			obj, err := m.SelectLoginToEdit()
			if err == nil {
				return NewLoginFormRedact(m.svc, obj), nil
			}
		case "t":
			m.objectType = "text"
			obj, err := m.SelectTextToEdit()
			if err == nil {
				return NewTextFormRedact(m.svc, obj), nil
			}
		case "c":
			m.objectType = "creditcard"
			obj, err := m.SelectCreditCardToEdit()
			if err == nil {
				return NewCardFormRedact(m.svc, obj), nil
			}
		case "f":
			m.objectType = "file"
			obj, err := m.SelectFileToEdit()
			if err == nil {
				return NewFileFormRedact(m.svc, obj), nil
			}
		}
	}

	return m, nil
}

func (m EditModel) View() string {
	var b string
	b = "Press [l] for Login, [t] for Text, [c] for Credit Card, [f] for File\n"
	return b
}

func (m EditModel) SelectCreditCardToEdit() (models.CreditCardData, error) {
	creditCards, err := m.svc.GetCardObjects()
	if err != nil {
		return models.CreditCardData{}, err
	}

	fmt.Println("Select a credit card to edit:")
	output := ""
	for i, card := range creditCards {
		output += fmt.Sprintf("[%d] Card ID: %s | Owner: %s || ", i, card.CardID, card.OwnerName)
	}
	fmt.Println(output)

	var selectedID string
	fmt.Print("Enter the ID of the card to edit: ")
	_, err = fmt.Scan(&selectedID)
	if err != nil {
		return models.CreditCardData{}, err
	}

	index, err := strconv.Atoi(selectedID)
	if err != nil || index < 0 || index >= len(creditCards) {
		return models.CreditCardData{}, err
	}

	selectedCard := creditCards[index]

	return selectedCard, nil
}

func (m EditModel) SelectLoginToEdit() (models.LoginData, error) {
	logins, err := m.svc.GetLoginObjects()
	if err != nil {
		return models.LoginData{}, err
	}

	fmt.Println("Select a login to edit:")
	var result string
	for i, login := range logins {
		result += fmt.Sprintf("[%d] Service: %s | Login: %s\n", i, login.ServiceName, login.Login)
	}
	fmt.Println(result)

	var selectedID string
	fmt.Print("Enter the ID of the card to edit: ")
	_, err = fmt.Scan(&selectedID)
	if err != nil {
		return models.LoginData{}, err
	}

	index, err := strconv.Atoi(selectedID)
	if err != nil || index < 0 || index >= len(logins) {
		return models.LoginData{}, err
	}

	selectedLogin := logins[index]

	return selectedLogin, nil
}

func (m EditModel) SelectTextToEdit() (models.TextData, error) {
	texts, err := m.svc.GetTextObjects()
	if err != nil {
		return models.TextData{}, err
	}

	fmt.Println("Select a credit card to edit:")
	var result string
	for i, text := range texts {
		result += fmt.Sprintf("[%d] Text: %s\n", i, text.Data)
	}
	fmt.Println(result)

	var selectedID string
	fmt.Print("Enter the ID of the card to edit: ")
	_, err = fmt.Scan(&selectedID)
	if err != nil {
		return models.TextData{}, err
	}

	index, err := strconv.Atoi(selectedID)
	if err != nil || index < 0 || index >= len(texts) {
		return models.TextData{}, err
	}

	selectedText := texts[index]

	return selectedText, nil
}

func (m EditModel) SelectFileToEdit() (models.FileData, error) {
	files, err := m.svc.GetFileObjects()
	if err != nil {
		return models.FileData{}, err
	}

	var result string
	for i, file := range files {
		result += fmt.Sprintf("[%d] FileName: %s\n", i, file.Path)
	}
	fmt.Println(result)

	var selectedID string
	fmt.Print("Enter the ID of the file to edit: ")
	_, err = fmt.Scan(&selectedID)
	if err != nil {
		return models.FileData{}, err
	}

	index, err := strconv.Atoi(selectedID)
	if err != nil || index < 0 || index >= len(files) {
		return models.FileData{}, err
	}

	selectedFile := files[index]

	filePath, err := openFileDialog()
	if err != nil {
		return models.FileData{}, err
	}
	selectedFile.Path = filePath

	return selectedFile, nil
}
