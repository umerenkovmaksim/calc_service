package app

import (
	"net/http"

	handler "github.com/umerenkovmaksim/calc_service/internal/handler"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", handler.CalcHandler)

	return mux
}
