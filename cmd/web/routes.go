package main

import (
	"github.com/ChrisCrawford1/Command/internal/handlers"
	internalMiddleware "github.com/ChrisCrawford1/Command/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

// Routes - Controls the main Chi Router
func Routes(requestHandler *handlers.RequestHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(internalMiddleware.ContentTypeMiddleware)

	router.Mount("/auth", AuthRouter(requestHandler))
	router.Mount("/users", UserRouter(requestHandler))

	return router
}

func AuthRouter(requestHandler *handlers.RequestHandler) http.Handler {
	authRouter := chi.NewRouter()
	authRouter.Post("/login", requestHandler.Login)
	return authRouter
}

func UserRouter(requestHandler *handlers.RequestHandler) http.Handler {
	userRouter := chi.NewRouter()
	userRouter.Group(func(r chi.Router) {
		r.Use(middleware2.ValidateJwtToken)
		r.Get("/me", requestHandler.GetMe)
	})
	return userRouter
}
