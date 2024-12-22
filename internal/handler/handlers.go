package handler

import (
	"encoding/json"
	"net/http"

	calc "github.com/umerenkovmaksim/calc_service/pkg/calculator"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func createResponse(w http.ResponseWriter, response CalcResponse, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		createResponse(w, CalcResponse{Error: "Invalid request method"}, http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var expression CalcRequest
	err := decoder.Decode(&expression)
	if err != nil {
		createResponse(w, CalcResponse{Error: "Invalid JSON format"}, http.StatusUnprocessableEntity)
		return
	}

	result, err := calc.Calc(expression.Expression)
	if err != nil {
		createResponse(w, CalcResponse{Error: err.Error()}, http.StatusUnprocessableEntity)
		return
	}
	createResponse(w, CalcResponse{Result: &result}, http.StatusOK)
}
