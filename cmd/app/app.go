package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ffajarpratama/boiler-api/config"
	"github.com/ffajarpratama/boiler-api/internal/http/handler"
	"github.com/ffajarpratama/boiler-api/internal/repository"
	"github.com/ffajarpratama/boiler-api/internal/usecase"
	"github.com/ffajarpratama/boiler-api/lib/mongo"
	"github.com/ffajarpratama/boiler-api/lib/postgres"
)

func Exec() error {
	cnf := config.New()

	pgClient, err := postgres.NewPostgresClient(cnf)
	if err != nil {
		return err
	}

	mongoClient, err := mongo.NewMongoClient(cnf)
	if err != nil {
		return err
	}

	repo := repository.New(pgClient, mongoClient)
	uc := usecase.New(cnf, repo, pgClient)
	r := handler.NewHTTPRouter(cnf, uc)

	addr := flag.String("http", fmt.Sprintf(":%d", cnf.App.Port), "HTTP listen address")
	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           r,
		ReadHeaderTimeout: 90 * time.Second,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Printf("[graceful-shutdown-time-out] \n%v\n", err.Error())
			}
		}()

		defer cancel()

		log.Println("graceful shutdown.....")

		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("[graceful-shutdown-error] \n%v\n", err.Error())
		}

		serverStopCtx()
	}()

	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("[http-server-failed] \n%v\n", err.Error())
		return err
	}

	<-serverCtx.Done()

	return nil
}
