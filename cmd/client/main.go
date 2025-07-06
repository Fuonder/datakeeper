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
	// login
	// add card data
	// get data
	// modify card data
	// get data

	// simulation of user input
	login := "testUser"
	password := "testPassword"
	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		panic(err)
	}

	user := models.User{
		ID:      0,
		Login:   login,
		PwdHash: password,
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
	//logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	//resp, err := client.R().
	//	SetHeader("Content-Type", "application/octet-stream").
	//	SetBody(string(cipherText)).
	//	Post("http://" + Flags.APIAddr.String() + "/register")
	//if err != nil {
	//	panic(err)
	//}
	//
	//if resp.StatusCode() != http.StatusOK {
	//	logger.Log.Debug("Error from server", zap.Any("message", string(resp.Body())))
	//	return
	//}
	//respBytes, err := cipherService.Decrypt(resp.Body())
	//if err != nil {
	//	panic(err)
	//}
	//
	//logger.Log.Debug("RESPONSE", zap.Any("response", string(respBytes)))
	//cookies := resp.Cookies()
	//authTokenR := ""
	//for _, cookie := range cookies {
	//	if cookie.Name == "auth_token" {
	//		authTokenR = cookie.Value
	//		logger.Log.Debug("Found auth token", zap.Any("cookie", authTokenR))
	//	}
	//}
	// login
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
