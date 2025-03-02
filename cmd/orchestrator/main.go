package main

import (
	"calc_service/internal/orchestrator"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions", orchestrator.GetExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionByIDHandler)
	http.HandleFunc("/internal/task", orchestrator.TaskHandler)

	log.Println("Starting orchestrator on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
