package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/auth"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/users"
	"net/http"
)

type Handlers struct {
	userSrv users.UserService
	authSrv auth.Service
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	logger.Log.Debug("GetLoginHandlerGet called")
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
//   - 401 Unauthorized: доступ для данного пользователя заблокирован.
//   - 500 Internal Server Error: внутренняя ошибка.
func (h Handlers) SaveCardHandlerPost(rw http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("SaveCardHandlerPost called")
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
	http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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
