package main

import (
	"github.com/ChrisCrawford1/Command/internal/handlers"
	internalMiddleware "github.com/ChrisCrawford1/Command/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
)

// Routes - Controls the main Chi Router
func Routes(requestHandler *handlers.RequestHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(internalMiddleware.ContentTypeMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	router.Mount("/auth", AuthRouter(requestHandler))
	router.Mount("/commands", CommandRouter(requestHandler))
	router.Mount("/users", UserRouter(requestHandler))

	return router
}

func AuthRouter(requestHandler *handlers.RequestHandler) http.Handler {
	authRouter := chi.NewRouter()
	authRouter.Post("/login", requestHandler.Login)
	return authRouter
}

func CommandRouter(requestHandler *handlers.RequestHandler) http.Handler {
	commandRouter := chi.NewRouter()
	commandRouter.Use(internalMiddleware.ValidateJwtToken)
	commandRouter.Post("/create", requestHandler.CreateCommand)
	commandRouter.Get("/all", requestHandler.GetAllCommands)
	commandRouter.Get("/{uuid}", requestHandler.GetCommand)
	commandRouter.Delete("/{uuid}", requestHandler.DeleteCommand)
	return commandRouter
}

func UserRouter(requestHandler *handlers.RequestHandler) http.Handler {
	userRouter := chi.NewRouter()
	userRouter.Use(internalMiddleware.ValidateJwtToken)
	userRouter.Get("/me", requestHandler.GetMe)
	return userRouter
}
