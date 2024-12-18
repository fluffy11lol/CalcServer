package main

import (
	"encoding/json"
	"fmt"
	"github.com/fluffy11lol/CalcServer/internal/config"
	"github.com/fluffy11lol/CalcServer/pkg/calculator"
	"log"
	"net/http"
)

type expression struct {
	Expression string `json:"expression"`
}
type result struct {
	Result float64 `json:"result"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	cfg := config.New()
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("starting server on", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Printf("failed to start server: %v", err)
		return
	}
}
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	expr := new(expression)
	err := json.NewDecoder(r.Body).Decode(expr)
	if err != nil {
		log.Printf("failed to decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		response := errorResponse{Error: err.Error()}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	res, err := calculator.Calc(expr.Expression)
	if err != nil {
		log.Printf("failed to calculate: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		response := errorResponse{Error: "expression is not valid: " + err.Error()}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	response := result{Result: res}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := errorResponse{Error: err.Error()}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	log.Printf("expression: %s, result: %f", expr.Expression, response.Result)
}
