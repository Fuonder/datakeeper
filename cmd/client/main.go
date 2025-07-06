package main

import (
	"encoding/json"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var (
	Token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RVc2VyIiwiZXhwIjoxNzUxODQ1MTI5fQ.CDzhwd5hxSWmD9Yq5gg6hgYedCHHYigwT-zD6Uom_l8"
	Login    = "testUser"
	Password = "testPassword"
)

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(Flags.LogLevel); err != nil {
		panic(fmt.Errorf("method main: %v", err))
	}
	logger.Log.Info("Flags parsed",
		zap.String("flags", Flags.String()))

	logger.Log.Info("Starting service")
	if err = run(); err != nil {
		logger.Log.Fatal("", zap.Error(err))
	}
}

func run() error {
	tester()
	return nil
}

func tester() {
	// register
	// RegisterTest()

	// login
	// LoginTest()

	// add card data
	// AddCardTest()

	// modify card data
	// UpdateCardTest()

	// get card by id
	GetCardTest()

	// get data

}

func AddCardTest() {
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}
	client := resty.New()

	card := models.CreditCardData{
		ID:        0,
		UserID:    0,
		CardID:    "1111-1111-1111-1111",
		OwnerName: "Test User",
		Metadata:  "Test credit card information",
	}

	//cipherCardID, err := cipherService.Encrypt([]byte(card.CardID))
	//if err != nil {
	//	panic(err)
	//}
	//card.CardID = string(cipherCardID)
	//cipherOwner, err := cipherService.Encrypt([]byte(card.OwnerName))
	//if err != nil {
	//	panic(err)
	//}
	//card.OwnerName = string(cipherOwner)

	jsonBytes, err := json.Marshal(card)
	if err != nil {
		panic(err)
	}
	cipherText, err := cipherService.Encrypt(jsonBytes)
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		SetCookie(&http.Cookie{Name: "auth_token", Value: Token}).
		Post("http://" + Flags.APIAddr.String() + "/data/card")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()),
			zap.Any("message", string(resp.Body())))
		return
	}
	respBytes, err := cipherService.Decrypt(resp.Body())
	if err != nil {
		panic(err)
	}

	updatedCardData := models.CreditCardData{}
	err = json.Unmarshal(respBytes, &updatedCardData)
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("Updated card", zap.Any("card", updatedCardData))
}

func UpdateCardTest() {
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}
	client := resty.New()

	card := models.CreditCardData{
		ID:        8,
		UserID:    0,
		CardID:    "1111-1111-1111-2222",
		OwnerName: "Test User2",
		Metadata:  "Test credit card information2",
	}

	//cipherCardID, err := cipherService.Encrypt([]byte(card.CardID))
	//if err != nil {
	//	panic(err)
	//}
	//card.CardID = string(cipherCardID)
	//cipherOwner, err := cipherService.Encrypt([]byte(card.OwnerName))
	//if err != nil {
	//	panic(err)
	//}
	//card.OwnerName = string(cipherOwner)

	jsonBytes, err := json.Marshal(card)
	if err != nil {
		panic(err)
	}
	cipherText, err := cipherService.Encrypt(jsonBytes)
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		SetCookie(&http.Cookie{Name: "auth_token", Value: Token}).
		Post("http://" + Flags.APIAddr.String() + "/data/card")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()),
			zap.Any("message", string(resp.Body())))
		return
	}
	respBytes, err := cipherService.Decrypt(resp.Body())
	if err != nil {
		panic(err)
	}

	updatedCardData := models.CreditCardData{}
	err = json.Unmarshal(respBytes, &updatedCardData)
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("Updated card", zap.Any("card", updatedCardData))
}

func GetCardTest() {
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}
	client := resty.New()

	resp, err := client.R().
		SetCookie(&http.Cookie{Name: "auth_token", Value: Token}).
		Get("http://" + Flags.APIAddr.String() + "/data/card/8")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()),
			zap.Any("message", string(resp.Body())))
		return
	}
	respBytes, err := cipherService.Decrypt(resp.Body())
	if err != nil {
		panic(err)
	}

	updatedCardData := models.CreditCardData{}
	err = json.Unmarshal(respBytes, &updatedCardData)
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("Updated card", zap.Any("card", updatedCardData))
}

func RegisterTest() {
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}

	user := models.User{
		ID:      0,
		Login:   Login,
		PwdHash: Password,
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	cipherText, err := cipherService.Encrypt(jsonBytes)
	if err != nil {
		panic(err)
	}

	client := resty.New()
	// registration
	logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		Post("http://" + Flags.APIAddr.String() + "/register")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server", zap.Any("message", string(resp.Body())))
		return
	}
	respBytes, err := cipherService.Decrypt(resp.Body())
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("RESPONSE", zap.Any("response", string(respBytes)))
	cookies := resp.Cookies()
	authTokenR := ""
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authTokenR = cookie.Value
			logger.Log.Debug("Found auth token", zap.Any("cookie", authTokenR))
		}
	}
}
func LoginTest() {
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}

	user := models.User{
		ID:      0,
		Login:   Login,
		PwdHash: Password,
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	cipherText, err := cipherService.Encrypt(jsonBytes)
	if err != nil {
		panic(err)
	}

	client := resty.New()
	logger.Log.Debug("NOW TRY TO LOGIN")
	logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		Post("http://" + Flags.APIAddr.String() + "/login")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server", zap.Any("message", string(resp.Body())))
		return
	}
	respBytes, err := cipherService.Decrypt(resp.Body())
	if err != nil {
		panic(err)
	}

	logger.Log.Debug("RESPONSE", zap.Any("response", string(respBytes)))
	cookies := resp.Cookies()
	authTokenL := ""
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			authTokenL = cookie.Value
			logger.Log.Debug("Found auth token", zap.Any("cookie", authTokenL))
		}
	}
}

// "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RVc2VyIiwiZXhwIjoxNzUxODQ1MTI5fQ.CDzhwd5hxSWmD9Yq5gg6hgYedCHHYigwT-zD6Uom_l8"
