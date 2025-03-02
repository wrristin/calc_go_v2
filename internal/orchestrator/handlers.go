package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add expression")

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Bad request:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Println("Expression received:", req.Expression)

	taskList := parseExpression(req.Expression)
	if taskList == nil {
		log.Println("Failed to parse expression:", req.Expression)
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	expressionID := taskList[0].ExpressionID
	log.Println("Generated expression ID:", expressionID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": expressionID})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received request to get a task")
		mutex.Lock()
		defer mutex.Unlock()

		for _, task := range tasks {
			if task.Status == "pending" {
				task.Status = "in progress"
				tasks[task.ID] = task
				log.Printf("Task assigned: ID=%s, ExpressionID=%s", task.ID, task.ExpressionID)
				json.NewEncoder(w).Encode(task)
				return
			}
		}
		log.Println("No tasks available")
		http.Error(w, "No tasks available", http.StatusNotFound)

	case http.MethodPost:
		log.Println("Received request to submit task result")
		var result struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			log.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		task, exists := tasks[result.ID]
		if !exists {
			log.Printf("Task not found: ID=%s", result.ID)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		log.Printf("Updating task: ID=%s, ExpressionID=%s", task.ID, task.ExpressionID)
		task.Result = result.Result
		task.Status = "completed"
		tasks[result.ID] = task

		// Обновляем статус выражения
		expression, exists := expressions[task.ExpressionID]
		if !exists {
			log.Printf("Expression not found: ID=%s", task.ExpressionID)
			return
		}

		allCompleted := true
		for _, t := range expression.Tasks {
			if tasks[t.ID].Status != "completed" {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			expression.Status = "completed"
			expression.Result = fmt.Sprintf("%.2f", result.Result)
			expressions[task.ExpressionID] = expression
			log.Printf("Expression completed: ID=%s, Result=%s", task.ExpressionID, expression.Result)
		} else {
			log.Printf("Not all tasks completed for ExpressionID=%s", task.ExpressionID)
		}

		w.WriteHeader(http.StatusOK)
	default:
		log.Println("Method not allowed:", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Listing all expressions")
	mutex.Lock()
	defer mutex.Unlock()

	var resp []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Result string `json:"result,omitempty"`
	}

	for id, expr := range expressions {
		resp = append(resp, struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Result string `json:"result,omitempty"`
		}{
			ID:     id,
			Status: expr.Status,
			Result: expr.Result,
		})
		log.Printf("Expression: ID=%s, Status=%s, Result=%s", id, expr.Status, expr.Result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": resp})
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	log.Println("Received request to get expression by ID:", id)
	mutex.Lock()
	defer mutex.Unlock()

	expression, exists := expressions[id]
	if !exists {
		log.Println("Expression not found:", id)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": expression,
	})
}
