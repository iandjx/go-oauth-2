package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/iandjx/go-oauth-2/pkg/auth"
	"github.com/iandjx/go-oauth-2/pkg/services/client"
	"github.com/iandjx/go-oauth-2/pkg/services/oauth"
	"github.com/iandjx/go-oauth-2/pkg/services/user"
)

const defaultPort = "8080"
const shutdownTimeout = 10 * time.Second

type Server struct {
	addr          string
	router        http.Handler
	OAuthService  oauth.Service
	UserServuce   user.Service
	ClientService client.Service
	AuthService   auth.AuthService
}

func New(addr string) *Server {
	s := new(Server)
	s.addr = addr
	return s
}
func (s *Server) setupRoutes() {

	userHandler := user.NewHandler(s.UserServuce)
	clientHandler := client.NewHandler(s.ClientService)

	r := mux.NewRouter()
	http.Handle("/", r)

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Handle("/register", userHandler.RegisterUser()).Methods(http.MethodPost)

	clientRouter := r.PathPrefix("/client").Subrouter()
	clientRouter.Use(auth.Middleware(s.AuthService))
	clientRouter.Handle("/create", clientHandler.CreateClient()).Methods(http.MethodPost)

}

func (s *Server) Run() error {
	s.setupRoutes()
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	return s.gracefulShutdown(srv, shutdownTimeout)
}

func (s *Server) gracefulShutdown(srv *http.Server, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		done <- srv.Shutdown(ctx)
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
	// return <-done
}
