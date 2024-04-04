package network

import (
	"Star13oy/medods/internal/service"
	"context"
	"net/http"
	"time"
)

// Server Структура для работы сервера (инстанс сервера, handler токенов для обработки http запросов)
type Server struct {
	srv          *http.Server
	tokenHandler *tokenHandler
}

// NewServer Конструктор сервера
func NewServer(tokenService *service.TokenService) *Server {
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Second * 5,
	}

	return &Server{
		srv:          srv,
		tokenHandler: newTokenHandler(tokenService),
	}
}

// Run Функция для начала работы сервера
func (s *Server) Run() error {
	//Первый маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре запроса
	http.HandleFunc("/get", s.tokenHandler.GetTokens)
	//Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
	http.HandleFunc("/refresh", s.tokenHandler.RefreshTokens)
	err := s.srv.ListenAndServe()
	return err
}

// Shutdown Функция для завершения работы сервера
func (s *Server) Shutdown() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.srv.Shutdown(ctx)
}
