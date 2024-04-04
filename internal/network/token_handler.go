package network

import (
	"Star13oy/medods/internal/service"
	"Star13oy/medods/internal/store"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Токен handler нужен для работы с http запросами и вызова необходимых функций tokenService
type tokenHandler struct {
	tokenService *service.TokenService
}

// Инициализация tokenHandler
func newTokenHandler(tokenService *service.TokenService) *tokenHandler {
	return &tokenHandler{
		tokenService: tokenService,
	}
}

// GetTokens GET запрос на получение пары Access и Refresh токенов для пользователя (в запросе также возвращается и GUID)
func (h *tokenHandler) GetTokens(w http.ResponseWriter, r *http.Request) {
	//Первый маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре запроса
	ctx := context.Background()

	// Достаем параметр из запроса
	guid := r.URL.Query().Get("GUID")

	// Создаем токен, рефреш токен из функции на основе GUID
	accessToken, refresh, err := h.tokenService.CreatePairTokens(ctx, guid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка формирования ключа: неверный guid :%v", err), 400)
		return
	}

	// Готовим все в нужный вид для вывода в Http ответе (вместе с guid, токеном и refresh токеном)
	u := store.User{GUID: guid, AccessToken: accessToken, RefreshTokens: refresh}
	resp, err := json.Marshal(u)
	if err != nil {
		http.Error(w, fmt.Sprintf("GetTokens: failed to marshal response: %v ", err), 400)
		return
	}
	// Отправляем ответ
	status, err := w.Write(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка отправки ответа %v, статус %v", err.Error(), status), 500)
		return
	}

	log.Printf("%v: успешная выдача токенов", status)

}

// RefreshTokens Get запрос на обновление пары Access и Refresh токенов для пользователя (в запросе также возвращается и GUID)
func (h *tokenHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	// Достаем из параметра refresh токен
	ctx := context.Background()

	refreshToken := r.URL.Query().Get("refresh_token")

	//* Берем у пользователя
	resultDecode, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		http.Error(w, "Неверный формат данных: предоставлен неправильный refresh токен", 400)
		return
	}

	var refreshTokenStruct store.RefreshToken               // Структура записи, содержащей hash Refresh токена в базе данных
	err = json.Unmarshal(resultDecode, &refreshTokenStruct) // преобразуем в структуру из пользовательской строки вида Base64
	if err != nil {
		http.Error(w, "Неверный формат данных: refresh токен имеет неправильную структуру", 400)
		return
	}
	//* Выгружаем данные из бд, ищем по GUID
	hashedRefresh, err := h.tokenService.FindRefreshByGuid(ctx, refreshTokenStruct.Guid) // Достаем системный refresh токен, чтобы сравнить в дальнейшем
	if err != nil {
		http.Error(w, "Неправильный refresh токен: не удалось найти пользователя", 400)
		return
	}

	// Проверяем токен сравнением с тем, что есть в базе данных на исходный GUID
	err = h.tokenService.CompareRefreshToken(refreshToken, hashedRefresh)
	if err != nil {
		http.Error(w, "Неверный формат данных: предоставлен неправильный refresh токен", 400)
		return
	}

	// Создаем пару токенов Access и Refresh
	accessToken, refresh, err := h.tokenService.CreatePairTokens(ctx, refreshTokenStruct.Guid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании ключей %v", err), 400)
		return
	}

	// Собираем ответ пользователю
	resp, err := json.Marshal(store.User{GUID: refreshTokenStruct.Guid, AccessToken: accessToken, RefreshTokens: refresh})
	if err != nil {
		http.Error(w, fmt.Sprintf("Ключи были повреждены или сформированы неправильно %v", err), 400)
		return
	}
	// Отправляем ответ
	status, err := w.Write(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка отправки результата %v, статус %v", err.Error(), status), 500)
		return
	}

	log.Println("Успешное обновление токенов")
}
