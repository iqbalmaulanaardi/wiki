package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"wiki/driver"
	handler "wiki/handler/http"
)
const (
	dbName = "go-mysql-crud"
	dbPass = "password"
	dbHost = "localhost"
	dbPort = "3306"
)
func main() {
	println("This is db", dbName, dbHost, dbPass, dbPort)

	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	pHandler := handler.NewPostHandler(connection)
	uHandler := handler.NewUserHandler(connection)


	r.Get("/posts", pHandler.Fetch)
	r.Get("/posts/{id}", pHandler.GetByID)
	r.Post("/posts", pHandler.Create)
	r.Put("/posts/{id}", pHandler.Update)
	r.Delete("/posts/{id}", pHandler.Delete)
	r.Get("/users", uHandler.Fetch)
	r.Get("/users/{id}", uHandler.GetById)
	r.Post("/signup", uHandler.Signup)
	r.Post("/signin", uHandler.Login)

	fmt.Println("Server listen at :8080")
	http.ListenAndServe(":8080", r)
}