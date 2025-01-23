package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"testcaseautomation/handlers"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("templates", tmpl)
	t, err := template.ParseFiles("templates/base.html", tmplPath)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTCHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /create-tc request")
	renderTemplate(w, "create_tc.html", map[string]interface{}{
		"Title": "Create Test Case",
	})
}

func manageTCHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /manage-tc request")
	renderTemplate(w, "manage_tc.html", map[string]interface{}{
		"Title": "Manage Test Cases",
	})
}

func apiLayerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /api-layer request")
	renderTemplate(w, "api_layer.html", map[string]interface{}{
		"Title": "API Layer",
	})
}

func main() {
	log.Println("Starting server on :8080")

	http.HandleFunc("/create-tc", createTCHandler)
	http.HandleFunc("/manage-tc", manageTCHandler)
	http.HandleFunc("/api-layer", apiLayerHandler)
	http.HandleFunc("/api/tcsteps", handlers.TCStepsHandler)
	http.HandleFunc("/api/tcsubmit", handlers.TCSubmitHandler)
	http.HandleFunc("/api/aitnumbers", handlers.AITNumbersHandler)
	http.HandleFunc("/api/expectedoutputs", handlers.ExpectedOutputsHandler)
	http.HandleFunc("/process-stats", handlers.GetProcessStats)
	http.HandleFunc("/resource-usage", handlers.GetResourceUsage)
	http.HandleFunc("/execute-command", handlers.ExecuteBackgroundCommandOnServer)
	http.HandleFunc("/collect-metrics", handlers.CollectMetrics)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
