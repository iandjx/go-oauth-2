package main

// TODO add handleError on all requests
// TODO add oauth handler
import (
	"fmt"
	"log"

	"github.com/iandjx/go-oauth-2/pkg/auth"
	"github.com/iandjx/go-oauth-2/pkg/config"
	"github.com/iandjx/go-oauth-2/pkg/services/client"
	"github.com/iandjx/go-oauth-2/pkg/services/user"
	"github.com/iandjx/go-oauth-2/pkg/sqlite"
	"github.com/iandjx/go-oauth-2/server"
)

type App struct {
	cfg     config.Config
	srv     *server.Server
	closeFn func() error
}

func newApp() *App {
	app := new(App)
	return app
}

func (a *App) loadConfig() error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	a.cfg = c
	return nil
}

func (a *App) setup() error {
	db, err := sqlite.New()
	if err != nil {
		log.Fatalln("Could not instantiate db")
	}
	userRepo := sqlite.NewUserRepository(db)
	userService := user.NewService(userRepo, a.cfg.JWTSecret)
	clientRepo := sqlite.NewClientRepository(db, a.cfg.JWTSecret)
	clientService := client.NewService(clientRepo, a.cfg.JWTSecret)

	authService := auth.NewService(a.cfg.JWTSecret)

	srv := server.New(a.cfg.Addr)
	srv.UserServuce = userService
	srv.ClientService = clientService
	srv.AuthService = authService

	a.srv = srv
	a.closeFn = func() error {
		if err = db.Close(); err != nil {
			log.Fatalln("could not close postgres:", err)
		}
		log.Println("app exited")

		return nil
	}
	return nil
}

func (a *App) run() error {
	fmt.Println("app started")
	return a.srv.Run()
}

func (a *App) close() error {
	if a.closeFn == nil {
		return nil
	}

	return a.closeFn()
}

func main() {

	app := newApp()
	defer app.close()

	if err := app.loadConfig(); err != nil {
		log.Fatalln("could not load config:", err)
	}

	if err := app.setup(); err != nil {
		log.Fatalln("error while setting up:", err)
	}

	if err := app.run(); err != nil {
		log.Fatalln("could not run app:", err)
	}
}
