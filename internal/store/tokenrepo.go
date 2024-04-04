package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// TokenRepo Интерфейс для TokenRepo
type TokenRepo interface {
	FindRefreshByGuid(ctx context.Context, guid string) (*UserMongo, error)
	UpsertRefreshToken(ctx context.Context, guid, hashedRefreshToken string) error
}

// Структура tokenRepo
type tokenRepo struct {
	client *mongo.Client // Mongodb клиент, который определяется в db
}

// NewTokenRepo конструктор TokenRepo
func NewTokenRepo(client *mongo.Client) TokenRepo {
	return &tokenRepo{
		client: client,
	}
}

// FindRefreshByGuid Функция для поиска Refresh токена внутри базы данных с помощью GUID
// Если найдена запись - вовзращаем её
func (r *tokenRepo) FindRefreshByGuid(ctx context.Context, guid string) (*UserMongo, error) {
	var user *UserMongo
	coll := r.client.Database("db").Collection("users")
	filter := bson.D{{Key: "Guid", Value: guid}} // Поиск записи guid и refresh токена (hash)
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, err
}

// UpsertRefreshToken Вставка, либо обновление новой записью текущей коллекции Refresh токенов
// Просходит после успешного выполнения Refresh токенов
func (r *tokenRepo) UpsertRefreshToken(ctx context.Context, guid, hashedRefreshToken string) error {
	// Формирование запроса вида Upsert (чтобы не плодить клонов пользователей с разными refresh токенами) внутри Mongodb
	filter := bson.D{
		{Key: "Guid", Value: guid}, // Ищем уникальные совпадения по Guid, тем самым обоюдно связываем Access токен и Refresh токен,
		// так как завязано все на пользователе и проверка идет в сравнении пользовательского refresh токена и системного.
		// Guid один и тот же, разные refresh токены не сойдутся из-за дальнейшей проверки (если пользователь поменял что-то).
	}
	// Ниже, вставляем данные в нужном формате и формируем запрос
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "Guid", Value: guid}, {Key: "BcryptHashRefresh", Value: hashedRefreshToken}}}}
	coll := r.client.Database("db").Collection("users")
	opts := options.Update().SetUpsert(true) // Опции, где включается Upsert
	resultUpdate, err := coll.UpdateOne(ctx, filter, update, opts)
	// Если все нормально - возвращаем refresh токен вида Base64 для пользователя
	if err != nil {
		return err
	}
	// Выводим результат в логи
	log.Default().Print(resultUpdate)

	return nil
}
