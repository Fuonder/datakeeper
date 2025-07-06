package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/auth"
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/dbservices"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"github.com/Fuonder/datakeeper.git/internal/objects/cards"
	"github.com/Fuonder/datakeeper.git/internal/objects/logins"
	"github.com/Fuonder/datakeeper.git/internal/objects/text"
	"github.com/Fuonder/datakeeper.git/internal/users"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Handlers struct {
	userSrv   users.UserService
	authSrv   auth.Service
	cipherSrv cipher.Service
	cardSrv   cards.Service
	loginSrv  logins.Service
	textSrv   text.Service
}

func NewHandlers(
	DBServices *dbservices.DatabaseServices,
	cryptoService cipher.Service) *Handlers {
	return &Handlers{
		userSrv:   DBServices.UserSrv,
		authSrv:   DBServices.AuthSrv,
		cipherSrv: cryptoService,
		cardSrv:   DBServices.CardSrv,
		loginSrv:  DBServices.LoginSrv,
		textSrv:   DBServices.TextSrv} // TODO: IMPLEMENT ME WHEN ALL SERVICES WILL BE DONE

}

//func NewHandlers(DBServices *dbservices.DatabaseServices) *Handlers {
//	return &Handlers{userSrv: DBServices.UserSrv,
//		walletSrv: DBServices.WalletSrv,
//		orderSrv:  DBServices.OrderSrv,
//		authSrv:   DBServices.AuthSrv}
//}

// RegisterHandlerPost регистрирует и аутентифицирует пользователя на сервере. Принимает данные в формате
// Content-Type: application/octet-stream, изначальный набор данных должен представлять собой объект типа models.User
//
// Возвращает:
//   - 200 OK: регистрация успешная, в ответе выставлен cookie в поле auth_token
//   - 400 Bad Request: в запросе переданы некорректные данные
//   - 409 Conflict: пользователь с заданным именем уже существует
//   - 500 Internal Server Error: внутренняя ошибка
func (h Handlers) RegisterHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("RegisterHandlerPost called")
	logger.Log.Debug("reading body")
	cipherText, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not read message"))
		return
	}
	logger.Log.Debug("Decrypting body")
	plainText, err := h.cipherSrv.Decrypt(cipherText)
	if err != nil {
		logger.Log.Debug("Can not decrypt message", zap.Error(err))
		SendResponse(rw, http.StatusBadRequest, []byte("Can not decrypt message"))
		return
	}
	logger.Log.Debug("Unmarshalling user object")
	userObject := models.User{}
	err = json.Unmarshal(plainText, &userObject)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not unmarshal message"))
		return
	}

	userObject.CreatedAt = time.Now()
	userObject.LastUpdate = time.Now()
	logger.Log.Debug("registering user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := h.authSrv.Register(ctx, userObject)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			SendResponse(rw, http.StatusConflict, []byte{})
			return
		}
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	logger.Log.Debug("register SUCCESS, setting cookie")
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	logger.Log.Debug("sending ok resp")
	respBody, err := h.cipherSrv.Encrypt([]byte("User created successfully"))
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
	}
	SendResponse(rw, http.StatusOK, respBody)
}

// LoginHandlerPost аутентифицирует пользователя на сервере. Принимает данные в формате
// // Content-Type: application/octet-stream, изначальный набор данных должен представлять собой объект типа models.User
//
// Возвращает:
//   - 200 OK: аутентификация успешная, в ответе выставлен cookie в поле auth_token
//   - 400 Bad Request: в запросе переданы некорректные данные
//   - 404 Not Found: пользователь с заданным именем не существует
//   - 500 Internal Server Error: внутренняя ошибка
func (h Handlers) LoginHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("LoginHandlerPost called")
	logger.Log.Debug("reading body")
	cipherText, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not read message"))
		return
	}
	logger.Log.Debug("Decrypting body")
	plainText, err := h.cipherSrv.Decrypt(cipherText)
	if err != nil {
		logger.Log.Debug("Can not decrypt message", zap.Error(err))
		SendResponse(rw, http.StatusBadRequest, []byte("Can not decrypt message"))
		return
	}
	logger.Log.Debug("Unmarshalling user object")
	userObject := models.User{}
	err = json.Unmarshal(plainText, &userObject)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not unmarshal message"))
		return
	}

	logger.Log.Debug("registering user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := h.authSrv.Login(ctx, userObject)
	if err != nil {
		logger.Log.Debug("error", zap.Error(err))
		if errors.Is(err, models.ErrWrongCredentials) {
			SendResponse(rw, http.StatusUnauthorized, []byte{})
			return
		}
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	logger.Log.Debug("register SUCCESS, setting cookie")
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	logger.Log.Debug("sending ok resp")
	respBody, err := h.cipherSrv.Encrypt([]byte("User login success"))
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
	}
	SendResponse(rw, http.StatusOK, respBody)
}

// DataHandlerGet возвращает список объектов, доступных пользователю на сервере.
// Для файлов не возвращает реальное наполнение файлов. Для получения непосредственного содержимого файла
// необходимо использовать GetFileHandlerGet. Возвращаемое значение имеет Content-Type: application/octet-stream.
// Возвращаемые данные представляют собой зашифрованный массив байт, созданный на основе данных типа models.ObjectList
//
// Возвращает:
//
//   - 200 OK: Объекты были найдены и возвращаются в ответе в виде зашифрованного массива байт.
//   - 204 No Content: для данного пользователя данные отсутствуют.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) DataHandlerGet(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("DataHandlerGet called")
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// SaveLoginHandlerPost используется для загрузки объекта типа models.LoginData на сервер. Принимает данные в
// зашифрованном виде. Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно добавлен. В ответе передается обновленный объект с установленным ID.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) SaveLoginHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("SaveLoginHandlerPost called")

	cipherText, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not read message"))
		return
	}

	logger.Log.Debug("Decrypting body")
	plainText, err := h.cipherSrv.Decrypt(cipherText)
	if err != nil {
		logger.Log.Debug("Can not decrypt message", zap.Error(err))
		SendResponse(rw, http.StatusBadRequest, []byte("Can not decrypt message"))
		return
	}

	loginObject := models.LoginData{}
	err = json.Unmarshal(plainText, &loginObject)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not unmarshal message"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusUnauthorized, []byte{})
		return
	}

	loginObject.UserID = UID
	loginObject.LastUpdate = time.Now()

	respLoginObject, err := h.loginSrv.AddNewLoginRecord(ctx, loginObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respJsonBytes, err := json.Marshal(respLoginObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	SendResponse(rw, http.StatusOK, respCipherText)
}

// GetLoginHandlerGet используется для скачивания объекта с сервера. Для объекта необходимо указать
// его ID в URL параметрах '/{object_id}'. Объект models.LoginData возвращается в закодированном, зашифрованном виде.
// Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно отправлен
//   - 400 Bad Request: Объект с заданным ID не найден.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) GetLoginHandlerGet(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("GetLoginHandlerGet called")

	loginIDString := chi.URLParam(r, "object_id")
	loginID, err := strconv.Atoi(loginIDString)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusUnauthorized, []byte{})
		return
	}

	loginObject, err := h.loginSrv.GetLoginRecord(ctx, loginID, UID)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respJsonBytes, err := json.Marshal(loginObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	SendResponse(rw, http.StatusOK, respCipherText)
}

// SaveTextHandlerPost используется для загрузки объекта типа models.TextData на сервер. Принимает данные в
// зашифрованном виде. Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно добавлен. В ответе передается обновленный объект с установленным ID.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) SaveTextHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("SaveTextHandlerPost called")

	cipherText, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not read message"))
		return
	}

	logger.Log.Debug("Decrypting body")
	plainText, err := h.cipherSrv.Decrypt(cipherText)
	if err != nil {
		logger.Log.Debug("Can not decrypt message", zap.Error(err))
		SendResponse(rw, http.StatusBadRequest, []byte("Can not decrypt message"))
		return
	}

	textObject := models.TextData{}
	err = json.Unmarshal(plainText, &textObject)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not unmarshal message"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusUnauthorized, []byte{})
		return
	}

	textObject.UserID = UID
	textObject.LastUpdate = time.Now()

	respTextObject, err := h.textSrv.AddNewTextRecord(ctx, textObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respJsonBytes, err := json.Marshal(respTextObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	SendResponse(rw, http.StatusOK, respCipherText)
}

// GetTextHandlerGet используется для скачивания объекта с сервера. Для объекта необходимо указать
// его ID в URL параметрах '/{object_id}'. Объект models.TextData возвращается в закодированном, зашифрованном виде.
// Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно отправлен
//   - 400 Bad Request: Объект с заданным ID не найден.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) GetTextHandlerGet(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("GetTextHandlerGet called")

	textIDString := chi.URLParam(r, "object_id")
	textID, err := strconv.Atoi(textIDString)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusUnauthorized, []byte{})
		return
	}

	textObject, err := h.textSrv.GetTextRecord(ctx, textID, UID)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respJsonBytes, err := json.Marshal(textObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	SendResponse(rw, http.StatusOK, respCipherText)
}

// SaveFileHandlerPost используется для загрузки объекта типа models.FileData на сервер. Принимает данные в
// зашифрованном виде. Content-Type: multipart/form-data.
//
// Возвращает:
//
//   - 200 OK: Объект успешно добавлен. В ответе передается объект models.FileData с установленным ID.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) SaveFileHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("SaveFileHandlerPost called")
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// GetFileHandlerGet используется для скачивания объекта с сервера. Для объекта необходимо указать
// его ID в URL параметрах '/{object_id}'. Объект models.FileData возвращается в закодированном, зашифрованном виде.
// Content-Type: multipart/form-data.
//
// Возвращает:
//
//   - 200 OK: Объект успешно отправлен
//   - 400 Bad Request: Объект с заданным ID не найден.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) GetFileHandlerGet(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("GetLoginHandlerGet called")
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

// SaveCardHandlerPost используется для загрузки объекта типа models.CreditCardData на сервер. Принимает данные в
// зашифрованном виде. Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно добавлен. В ответе передается обновленный объект с установленным ID.
//   - 400 Bad Request: некорректный формат сообщения
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) SaveCardHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("SaveCardHandlerPost called")
	logger.Log.Debug("reading body")
	cipherText, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not read message"))
		return
	}
	logger.Log.Debug("Decrypting body")
	plainText, err := h.cipherSrv.Decrypt(cipherText)
	if err != nil {
		logger.Log.Debug("Can not decrypt message", zap.Error(err))
		SendResponse(rw, http.StatusBadRequest, []byte("Can not decrypt message"))
		return
	}
	logger.Log.Debug("Unmarshalling user object")
	cardObject := models.CreditCardData{}
	err = json.Unmarshal(plainText, &cardObject)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte("Can not unmarshal message"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	cardObject.UserID = UID
	cardObject.LastUpdate = time.Now()

	respCardObject, err := h.cardSrv.AddNewCardRecord(ctx, cardObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	respJsonBytes, err := json.Marshal(respCardObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	SendResponse(rw, http.StatusOK, respCipherText)
}

// GetCardHandlerGet используется для скачивания объекта с сервера. Для объекта необходимо указать
// его ID в URL параметрах '/{object_id}'. Объект models.CreditCardData возвращается в закодированном, зашифрованном виде.
// Content-Type: application/octet-stream.
//
// Возвращает:
//
//   - 200 OK: Объект успешно отправлен
//   - 400 Bad Request: Объект с заданным ID не найден.
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) GetCardHandlerGet(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("GetLoginHandlerGet called")

	cardIDString := chi.URLParam(r, "object_id")
	cardID, err := strconv.Atoi(cardIDString)
	if err != nil {
		SendResponse(rw, http.StatusBadRequest, []byte{})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UID, err := h.getUserID(ctx, r)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}

	cardObject, err := h.cardSrv.GetCardRecord(ctx, cardID, UID)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	respJsonBytes, err := json.Marshal(cardObject)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	respCipherText, err := h.cipherSrv.Encrypt(respJsonBytes)
	if err != nil {
		SendResponse(rw, http.StatusInternalServerError, []byte{})
		return
	}
	SendResponse(rw, http.StatusOK, respCipherText)
}

// SendResponse вспомогательная функция, предназначенная для формирования http ответа и его отправки.
// Ответ записывается в заданный rw http.ResponseWriter. Ответ будет содержать заданный status и сообщение message.
func SendResponse(rw http.ResponseWriter, status int, message []byte) {
	rw.WriteHeader(status)
	if len(message) == 0 {
		_, _ = rw.Write([]byte(http.StatusText(status)))
		return
	}
	_, _ = rw.Write(message)
}

// getUserID внутренняя функция, предназначенная для получения ID пользователя из базы на основе аутентификационного
// токена.
func (h Handlers) getUserID(ctx context.Context, r *http.Request) (int, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// If the cookie is not found, return an error
			return 0, fmt.Errorf("cookie not found")
		}
		return 0, fmt.Errorf("error retrieving cookie: %v", err)
	}
	tokenString := cookie.Value
	UID, err := h.authSrv.GetUIDFromJWT(ctx, tokenString)
	if err != nil {
		return 0, fmt.Errorf("error retrieving user ID: %v", err)
	}
	return UID, nil
}
