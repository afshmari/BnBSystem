package main

import (
	"Hotel/pkg/config"
	"Hotel/pkg/handlers"
	"Hotel/pkg/render"
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("This port is on port %s", portNumber))

	err = http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
