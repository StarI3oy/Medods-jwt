package config

const (
	MONGODB_HOST = "mongodb://localhost:27017" // адрес MongoDB (поменять на нижний, если будет compose)
	//MONGODB_HOST = "mongodb://mongo:27017"
	SECRET_KEY = "VERY_SECRET_KEY" // секретный ключ для генерации JWT
)
