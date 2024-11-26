package handlers

import (
	"fmt"
	"net/http"
)

type TCStep struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func TCStepsHandler(w http.ResponseWriter, r *http.Request) {
	tcsteps := []TCStep{
		{ID: "1", Value: "Step1"},
		{ID: "2", Value: "Step2"},
		{ID: "3", Value: "Step3"},
	}

	w.Header().Set("Content-Type", "text/html")

	for _, step := range tcsteps {
		fmt.Fprintf(w, `<option value="%s">%s</option>`, step.ID, step.Value)
	}
}
