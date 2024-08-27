package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"weex_admin/cmd/http/di"
	"weex_admin/internal/context"
	"weex_admin/internal/infrastructure/database"
	"weex_admin/internal/infrastructure/jwt"
)

func main() {
	var (
		doneChannel  = make(chan struct{})
		errorChannel = make(chan error)

		gormConf = &database.DatabaseConfig{
			Host:     "127.0.0.1",
			User:     "root",
			Port:     3306,
			Password: "rootpassword",
			DbName:   "demo",
			DbType:   "mysql",
		}

		jwtConf = &jwt.JWTManagerConf{
			Secret: []byte("mysecret"),
		}

		appContext = &context.AppContext{
			ErrorChannel: errorChannel,
			DoneChannel:  doneChannel,
		}
	)

	// Setup Gin router
	app, err := di.InitializeApp(
		appContext,
		gormConf,
		jwtConf,
	)
	if err != nil {
		log.Fatalf("failed to initialize engine: %v", err)
	}

	// Seed database
	app.DatabaseInitializer.SeedDatabase()

	endChannel := deferEnd(app)

	// Run HTTP server
	if err := app.Engine.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	close(doneChannel)

	select {
	case <-endChannel:
		log.Println("server shutdown")
	case <-time.After(5 * time.Minute):
		log.Println("server shutdown timeout")
	}
}

func deferEnd(app *di.App) <-chan struct{} {
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		wg.Add(1)
		<-app.AppContext.DoneChannel

		if sqlDb, err := app.DB.DB(); err != nil {
			app.AppContext.ErrorChannel <- fmt.Errorf("failed to get database: %v", err)
		} else {
			if err := sqlDb.Close(); err != nil {
				app.AppContext.ErrorChannel <- fmt.Errorf("failed to close database: %v", err)
			}
		}
	}()

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}
