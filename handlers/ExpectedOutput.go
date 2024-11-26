package handlers

import (
	"fmt"
	"net/http"
)

type ExpectedOutput struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func ExpectedOutputsHandler(w http.ResponseWriter, r *http.Request) {
	expectedOutputs := []ExpectedOutput{
		{ID: "1", Value: "Output-1"},
		{ID: "2", Value: "Output-2"},
		{ID: "3", Value: "Output-3"},
	}

	w.Header().Set("Content-Type", "text/html")

	for _, output := range expectedOutputs {
		fmt.Fprintf(w, `<option value="%s">%s</option>`, output.ID, output.Value)
	}
}
