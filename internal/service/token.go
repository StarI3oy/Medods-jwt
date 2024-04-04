package service

import (
	"Star13oy/medods/internal/config"
	"Star13oy/medods/internal/store"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// TokenService отвечает за основные функции для работы с токенами
type TokenService struct {
	repo store.TokenRepo // Некоторые функции требуют использование базы данных MongoDB и для этого была настроена Repo
}

// NewTokenService конструктор TokenService
func NewTokenService(repo store.TokenRepo) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

// CreateRefreshToken Создание Refresh токена
func (s *TokenService) CreateRefreshToken(ctx context.Context, guid string, timePoint time.Time) (string, error) {
	// Создаем структуру для хранения refresh ключей в бд
	rtjStruct := store.RefreshToken{Guid: guid, Expires: timePoint}

	// Собираем в байты, потом в строку
	refTokenJson, err := json.Marshal(rtjStruct)
	if err != nil {
		return "", err
	}
	refTokenString := string(refTokenJson)

	// Далее - этап формирования двух видов refresh токена
	// Первый (у пользователя) - сформирован в виде Base64 (ТЗ), основываясь на строковом представлении структуры рефреш токена
	// Второй (в бд) - сформирован в виде хэша (ТЗ) с использованием bcrypt для хранения в бд

	//* Передается пользователю
	refreshToken := base64.StdEncoding.EncodeToString([]byte(refTokenString))

	//* Передается в бд
	bytes, err := bcrypt.GenerateFromPassword(
		refTokenJson,
		8,
	)
	if err != nil {
		return "", err
	}
	hashedRefreshToken := string(bytes)

	err = s.repo.UpsertRefreshToken(ctx, guid, hashedRefreshToken)
	if err != nil {
		return "", err
	}

	return refreshToken, nil

}

// CompareRefreshToken декодирует пользовательский Refresh токен вида base64 и сравнивает с Hash который хранится в базе данных
func (s *TokenService) CompareRefreshToken(refreshToken, hash string) error {
	// Проверяем на совпадение refresh токена системы и пользователя (для этого преобразовываем пользовательский в исходную строку)
	resultDecode, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return err
	}
	// Далее - сравниваем hash и строку refresh токена (вернется nil если хэш был образован от исходной, т.е. refresh токен не поменял пользователь)
	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		resultDecode,
	)
	return err
}

// CreatePairTokens Создание пары токенов Access и Refresh
func (s *TokenService) CreatePairTokens(
	ctx context.Context,
	guid string,
) (accessToken, refreshToken string, err error) {
	if guid == "" {
		return "", "", errors.New("guid не предоставлен")
	}
	timeMoment := time.Now()
	// Создаем access токен
	accessToken, err = s.createToken(guid, timeMoment)
	if err != nil {
		return "", "", err
	}

	// Создаем refresh токен
	refreshToken, err = s.CreateRefreshToken(ctx, guid, timeMoment)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// FindRefreshByGuid Поиск hash Refresh токена, который хранится в базе данных (обращаемся к функции внутри repo)
func (s *TokenService) FindRefreshByGuid(ctx context.Context, guid string) (string, error) {
	user, err := s.repo.FindRefreshByGuid(ctx, guid)
	if err != nil {
		return "", err
	}

	return user.BcryptHashRefresh, nil
}

// createToken Создание Access токена
func (s *TokenService) createToken(guid string, timePoint time.Time) (string, error) {
	// Создаем токен на основе JWT через функцию
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, // Используем для создания JWT SHA512, по тз
		jwt.MapClaims{ // Наполнение не было написано в ТЗ, поэтому своё - GUID и expires
			"GUID":    []byte(guid),
			"expires": timePoint.AddDate(0, 10, 0),
		})
	res, err := t.SignedString([]byte(config.SECRET_KEY))
	return res, err

}
