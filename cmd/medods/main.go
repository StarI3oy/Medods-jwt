package main

import (
	"Star13oy/medods/internal/db"
	"Star13oy/medods/internal/network"
	"Star13oy/medods/internal/service"
	"Star13oy/medods/internal/store"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	// Подключаемся к MongoDB по подготовленной функции с настройками
	client, err := db.ConnectMongo()
	if err != nil {
		log.Panicf("failed to connect mongo: %v", err)
	}
	// После всех операций, необходимо будет отключиться от Mongodb
	defer client.Disconnect(context.Background())

	tokenRepo := store.NewTokenRepo(client) // tokenRepo содержит ( и был создан) для взаимодействия с базой данных
	// Поэтому передаем внутрь client и пользуемся функциями через сервис

	tokenService := service.NewTokenService(tokenRepo) // Инициализируем сервис поверх repo для работы с базой данных
	// Внутри все необходимые функции для работы с токенами, но взаимодействие с базой данных проходит через repo

	srv := network.NewServer(tokenService) // Сервис для контроля сервера
	// Внутри выведены две основные функции Run и Shutdown

	go func() { // Запускаем инстанс сервера через горутин
		err = srv.Run()
		if err != nil {
			log.Panicf("Сервер упал: %v", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	srv.Shutdown()

}
