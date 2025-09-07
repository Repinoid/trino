package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"triner/internal/config"
	"triner/internal/dbase"
	"triner/internal/handlera"
	"triner/internal/models"

	"github.com/gorilla/mux"

	_ "github.com/trinodb/trino-go-client/trino"
)

func main() {

	ctx := context.Background()

	// уровень логирования по умолчанию Info

	if err := Run(ctx); err != nil {
		models.Logger.Error(err.Error())
	}

}

func Run(ctx context.Context) (err error) {

	Level := slog.LevelInfo

	// Если есть флаг -debug
	debugFlag := flag.Bool("debug", false, "установка Минимального уровня логирования DEBUG")
	flag.Parse()
	if *debugFlag {
		Level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     Level,
		AddSource: true, // Добавлять информацию об исходном коде
	})
	models.Logger = slog.New(handler)
	slog.SetDefault(models.Logger)
	models.Logger.Debug("Log", "level", Level)

	postgres, err := dbase.NewPostgresPool(context.Background(), models.DSN)
	if err != nil {
		log.Fatalln("NewPostgresPool", "fault", err)
		return
	}
	defer postgres.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", handlera.DBPinger).Methods("GET")
	router.HandleFunc("/t", handlera.TrinoPinger).Methods("GET")

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Используем AppHost (или 0.0.0.0) и AppPort для HTTP-сервера
	serverAddr := fmt.Sprintf("%s:%d", config.Configuration.AppHost, config.Configuration.AppPort)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		fmt.Printf("\nServer started on %s\n\n", serverAddr)
		models.Logger.Info("Server started", "on", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидаем SIGINT (Ctrl+C) или SIGTERM
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel() // При получении сигнала отменяем контекст
	}()

	// Блокируемся, пока контекст не отменён
	<-ctx.Done()

	// Graceful shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		models.Logger.Error("Shutdown", "error", err.Error())
	} else {
		models.Logger.Info("Server stopped gracefully")
	}

	return

}
