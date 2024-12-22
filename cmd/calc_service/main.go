package main

import (
	"net/http"

	app "github.com/umerenkovmaksim/calc_service/internal/app"
)

func main() {
	router := app.NewRouter()
	http.ListenAndServe(":8080", router)
}
