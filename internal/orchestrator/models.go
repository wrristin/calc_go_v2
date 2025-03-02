package orchestrator

import (
	"sync"
	"time"
)

var (
	expressions = make(map[string]Expression)
	tasks       = make(map[string]Task)
	mutex       = &sync.Mutex{}
)

type Expression struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
	Tasks  []Task `json:"-"`
}

type Task struct {
	ID            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
	Status        string        `json:"status"`
	Result        float64       `json:"result,omitempty"`
	ExpressionID  string        `json:"expression_id"`
}
