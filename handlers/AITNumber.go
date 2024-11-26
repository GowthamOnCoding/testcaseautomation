package handlers

import (
	"fmt"
	"net/http"
)

type AITNumber struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func AITNumbersHandler(w http.ResponseWriter, r *http.Request) {
	aitnumbers := []AITNumber{
		{ID: "1", Value: "AIT-1001"},
		{ID: "2", Value: "AIT-1002"},
		{ID: "3", Value: "AIT-1003"},
	}

	w.Header().Set("Content-Type", "text/html")

	for _, ait := range aitnumbers {
		fmt.Fprintf(w, `<option value="%s">%s</option>`, ait.ID, ait.Value)
	}
}
