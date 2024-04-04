package store

import (
	"time"
)

// User Структура возвращаемого запроса
type User struct {
	GUID          string
	AccessToken   string
	RefreshTokens string
}

// UserMongo Структура для хранения в базе данных
type UserMongo struct {
	Guid              string
	BcryptHashRefresh string
}

// RefreshToken Структура Refresh токена, который хранится в базе данных
type RefreshToken struct {
	Guid    string
	Expires time.Time
}
