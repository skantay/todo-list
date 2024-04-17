package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/skantay/todo-list/config"
	v1 "github.com/skantay/todo-list/internal/controller/http/v1"
	"github.com/skantay/todo-list/internal/repository"
	"github.com/skantay/todo-list/internal/usecase"
	"github.com/skantay/todo-list/pkg/httpserver"
	"github.com/skantay/todo-list/pkg/log"
	"github.com/skantay/todo-list/pkg/mongodb"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	configPath = "config/config.yaml"
)

func Run() error {
	// Загрузка конфиг параметров
	cfg, err := config.New(configPath)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	// Указываем параметры для установки соединения с MongoDB
	options := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%s",
			cfg.MongoDB.User,
			cfg.MongoDB.Password,
			cfg.MongoDB.Host,
			cfg.MongoDB.Port,
		),
	)

	// Прокидываем контекст в MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), mongodb.DefaultTimeout)
	defer cancel()

	// Установка соединения с MongoDB
	client, err := mongodb.Connect(ctx, options)
	if err != nil {
		return fmt.Errorf("error connecting to mongodb: %w", err)
	}

	// Здесь указываются с какими коллекциями будет репозиторный слой работать
	collections := repository.Collections{
		Task: "task",
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}
	// Инициализация логгера
	logger := log.InitSlog(opts)

	// иньекций зависимостей
	repository := repository.New(client, "taskdb", collections, logger)
	usecase := usecase.New(repository, logger)

	router := gin.Default()
	v1.Set(router, usecase, logger)

	logger.Info("starting server on", "host", cfg.Server.Host, "port", cfg.Server.Port)

	// Запуск сервера
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.Server.Port),
	)

	// gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнал с двух каналов
	// канал interrupt - с него ожидаем syscall.SIGTERN или же просто "CTRL+C"
	// функция httpServer.Notify() возвращает канал, и с этого канала ожидаем какие либо ошибки при запуске server.ListenAndServe()
	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		slog.Error("app - Run - httpServer.Notify: %w", err)
	}

	slog.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("app - Run - httpServer.Shutdown: %w", err)
	}

	return err
}
