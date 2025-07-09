package cliservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/client/storage"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Service struct {
	st         *storage.InMemoryStorage
	cp         cipher.Encryptor
	remoteAddr string
	ticker     *time.Ticker
	client     *resty.Client
}

func NewService(aesKey []byte, remoteAddr string, pollInterval time.Duration) (*Service, error) {
	var err error

	app := new(Service)
	app.st = storage.NewInMemoryStorage()

	app.cp, err = cipher.NewAES256Cipher(aesKey)
	if err != nil {
		return nil, err
	}

	app.remoteAddr = remoteAddr
	app.ticker = time.NewTicker(pollInterval)
	app.client = resty.New()
	go app.Poll(context.TODO())
	return app, nil
}

func (s *Service) Register(user models.User) error {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	cipherText, err := s.cp.Encrypt(jsonBytes)
	if err != nil {
		return err
	}
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		Post("http://" + s.remoteAddr + "/register")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server")
		return fmt.Errorf(string(resp.Body()))
	}
	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return err
	}

	logger.Log.Debug("RESPONSE", zap.Any("response", string(respBytes)))
	cookies := resp.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			logger.Log.Debug("Found auth token", zap.Any("cookie", cookie.Value))
			s.st.SetToken(cookie.Value)
		}
	}
	return nil
}

func (s *Service) Login(user models.User) error {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	cipherText, err := s.cp.Encrypt(jsonBytes)
	if err != nil {
		return err
	}

	logger.Log.Debug("NOW TRY TO LOGIN")
	logger.Log.Debug("Sending body of size:", zap.Any("size", len(cipherText)))
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText)).
		Post("http://" + s.remoteAddr + "/login")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server")
		return fmt.Errorf(string(resp.Body()))
	}
	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return err
	}

	logger.Log.Debug("RESPONSE", zap.Any("response", string(respBytes)))
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			logger.Log.Debug("Found auth token", zap.Any("cookie", cookie.Value))
			s.st.SetToken(cookie.Value)
		}
	}
	return nil
}

func (s *Service) Poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Debug("client poller stopped")
			return
		case <-s.ticker.C:
			if s.st.GetToken() == "" {
				//logger.Log.Debug("Not authenticated yet, skipping polling")
				continue
			}
			data, err := s.FetchData()
			if err != nil {
				//logger.Log.Warn("poller fetch failed: ", zap.Error(err))
				continue
			}

			for _, item := range data.TextObjects {
				_ = s.st.AddItem(item)
			}
			for _, item := range data.LoginObjects {
				_ = s.st.AddItem(item)
			}
			for _, item := range data.CreditCardObjects {
				_ = s.st.AddItem(item)
			}
			for _, item := range data.FileObjects {
				_ = s.st.AddItem(item)
			}
		}
	}
}

func (s *Service) FetchData() (models.ObjectList, error) {

	resp, err := s.client.R().
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Get(fmt.Sprintf("http://%s/data", s.remoteAddr))
	if err != nil {
		return models.ObjectList{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()))
		return models.ObjectList{}, fmt.Errorf(string(resp.Body()))
	}

	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return models.ObjectList{}, err
	}

	var objList models.ObjectList
	if err := json.Unmarshal(respBytes, &objList); err != nil {
		return models.ObjectList{}, err
	}
	//logger.Log.Debug("Fetched data list")
	return objList, nil
}

func (s *Service) AddItem(item interface{}) error {
	request, resultObject, endpoint, err := s.createRequestForItem(item)
	if err != nil {
		return err
	}

	resp, err := request.
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Post(fmt.Sprintf("http://%s%s", s.remoteAddr, endpoint))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server", zap.Int("Status", resp.StatusCode()))
		return fmt.Errorf(string(resp.Body()))
	}

	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBytes, resultObject); err != nil {
		return err
	}

	if err := s.st.AddItem(resultObject); err != nil {
		return err
	}

	logger.Log.Debug("Updated object", zap.Any("value", resultObject))
	return nil
}

func (s *Service) GetTextItem(idx int) (models.TextData, error) {
	return s.st.GetTextItemByIndex(idx)
}
func (s *Service) GetCardItem(idx int) (models.CreditCardData, error) {
	return s.st.GetCardItemByIndex(idx)
}
func (s *Service) GetLoginItem(idx int) (models.LoginData, error) {
	return s.st.GetLoginItemByIndex(idx)
}
func (s *Service) GetFileItem(idx int) (models.FileData, error) {
	return s.st.GetFileItemByIndex(idx)
}

func (s *Service) GetTextObjects() ([]models.TextData, error) {
	return s.st.GetTextObjects()
}
func (s *Service) GetCardObjects() ([]models.CreditCardData, error) {
	return s.st.GetCardObjects()
}
func (s *Service) GetLoginObjects() ([]models.LoginData, error) {
	return s.st.GetLoginObjects()
}
func (s *Service) GetFileObjects() ([]models.FileData, error) {
	return s.st.GetFileObjects()
}

func (s *Service) FetchFileItem(item models.FileData) (FileItem models.FileData, FileContents []byte, error error) {
	resp, err := s.client.R().
		SetDoNotParseResponse(true).
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Get(fmt.Sprintf("http://%s/data/file/%d", s.remoteAddr, item.ID))
	if err != nil {
		return models.FileData{}, []byte{}, err
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != http.StatusOK {
		body, _ := io.ReadAll(resp.RawBody())
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()),
			zap.Any("message", string(body)))
		return
	}

	// Читаем multipart тело
	contentType := resp.Header().Get("Content-Type")
	if contentType == "" {
		return models.FileData{}, []byte{}, fmt.Errorf("missing Content-Type")
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return models.FileData{}, []byte{}, fmt.Errorf("invalid content type: not multipart")
	}

	mr := multipart.NewReader(resp.RawBody(), params["boundary"])

	var fileMeta models.FileData
	var fileContentsEncrypted []byte

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return models.FileData{}, []byte{}, err
		}

		switch part.FormName() {
		case "meta":
			metaBytes, err := io.ReadAll(part)
			if err != nil {
				return models.FileData{}, []byte{}, err
			}
			// 🔐 Расшифровываем только часть meta
			plainMeta, err := s.cp.Decrypt(metaBytes)
			if err != nil {
				return models.FileData{}, []byte{}, err
			}
			if err := json.Unmarshal(plainMeta, &fileMeta); err != nil {
				return models.FileData{}, []byte{}, err
			}

		case "file":
			fileContentsEncrypted, err = io.ReadAll(part)
			if err != nil {
				return models.FileData{}, nil, fmt.Errorf("read file: %w", err)
			}
			//outPath := "./downloaded_" + filepath.Base(part.FileName())
			//out, err := os.Create(outPath)
			//if err != nil {
			//	return models.FileData{}, []byte{}, err
			//}
			//defer out.Close()
			//
			//if _, err := io.Copy(out, part); err != nil {
			//	return models.FileData{}, []byte{}, err
			//}
			//logger.Log.Debug("File saved", zap.String("path", outPath))
		}
	}
	if fileContentsEncrypted == nil {
		return models.FileData{}, nil, fmt.Errorf("file contents not found in response")
	}

	decryptedFileContents, err := s.cp.Decrypt(fileContentsEncrypted)
	if err != nil {
		return models.FileData{}, nil, fmt.Errorf("decrypt file: %w", err)
	}

	return fileMeta, decryptedFileContents, nil

}

func (s *Service) FetchTextItem(item models.TextData) (TextItem models.TextData, error error) {
	resp, err := s.client.R().
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Get(fmt.Sprintf("http://%s/data/text/%d", s.remoteAddr, item.ID))
	if err != nil {
		return models.TextData{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()))
		return models.TextData{}, fmt.Errorf(string(resp.Body()))
	}

	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return models.TextData{}, err
	}

	var text models.TextData
	if err := json.Unmarshal(respBytes, &text); err != nil {
		return models.TextData{}, err
	}
	logger.Log.Debug("Fetched text", zap.Any("text", text))
	if err := s.st.AddItem(text); err != nil {
		return models.TextData{}, err
	}
	return text, nil
}

func (s *Service) FetchLoginItem(item models.LoginData) (LoginItem models.LoginData, error error) {
	resp, err := s.client.R().
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Get(fmt.Sprintf("http://%s/data/login/%d", s.remoteAddr, item.ID))
	if err != nil {
		return models.LoginData{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()))
		return models.LoginData{}, fmt.Errorf(string(resp.Body()))
	}

	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return models.LoginData{}, err
	}

	var login models.LoginData
	if err := json.Unmarshal(respBytes, &login); err != nil {
		return models.LoginData{}, err
	}
	logger.Log.Debug("Fetched login", zap.Any("login", login))
	if err := s.st.AddItem(login); err != nil {
		return models.LoginData{}, err
	}
	return login, nil
}

func (s *Service) FetchCardItem(item models.CreditCardData) (CardItem models.CreditCardData, error error) {
	resp, err := s.client.R().
		SetCookie(&http.Cookie{Name: "auth_token", Value: s.st.GetToken()}).
		Get(fmt.Sprintf("http://%s/data/card/%d", s.remoteAddr, item.ID))
	if err != nil {
		return models.CreditCardData{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Debug("Error from server",
			zap.Any("Status", resp.StatusCode()))
		return models.CreditCardData{}, fmt.Errorf(string(resp.Body()))
	}
	respBytes, err := s.cp.Decrypt(resp.Body())
	if err != nil {
		return models.CreditCardData{}, err
	}

	updatedCardData := models.CreditCardData{}
	err = json.Unmarshal(respBytes, &updatedCardData)
	if err != nil {
		return models.CreditCardData{}, err
	}

	logger.Log.Debug("Fetched card", zap.Any("card", updatedCardData))

	if err := s.st.AddItem(updatedCardData); err != nil {
		return models.CreditCardData{}, err
	}
	return updatedCardData, nil
}

func (s *Service) prepareEncryptedRequest(endpoint string, obj interface{}, result interface{}) (*resty.Request, string, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, "", err
	}
	cipherText, err := s.cp.Encrypt(jsonBytes)
	if err != nil {
		return nil, "", err
	}

	request := s.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(string(cipherText))

	return request, endpoint, nil
}
func (s *Service) createRequestForItem(item interface{}) (*resty.Request, interface{}, string, error) {
	var (
		request      *resty.Request
		resultObject interface{}
		endpoint     string
		err          error
	)
	switch v := item.(type) {
	case models.TextData:
		endpoint = "/data/text"
		resultObject = &models.TextData{}
		request, _, err = s.prepareEncryptedRequest(endpoint, v, resultObject)

	case models.LoginData:
		endpoint = "/data/login"
		resultObject = &models.LoginData{}
		request, _, err = s.prepareEncryptedRequest(endpoint, v, resultObject)

	case models.CreditCardData:
		endpoint = "/data/card"
		resultObject = &models.CreditCardData{}
		request, _, err = s.prepareEncryptedRequest(endpoint, v, resultObject)

	case models.FileData:
		endpoint = "/data/file"
		resultObject = &models.FileData{}

		metaBytes, err := json.Marshal(v)
		if err != nil {
			return nil, nil, "", err
		}
		cipherText, err := s.cp.Encrypt(metaBytes)
		if err != nil {
			return nil, nil, "", err
		}

		request = s.client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetFile("file", v.Path).
			SetFormData(map[string]string{
				"meta": string(cipherText),
			})

	default:
		return nil, nil, "", fmt.Errorf("unknown type %v", reflect.TypeOf(item))
	}

	if err != nil {
		return nil, nil, "", err
	}

	return request, resultObject, endpoint, nil
}

func (s *Service) Close() {
	s.ticker.Stop()
}
