package main

import (
	"net/http"
	"testcaseautomation/handlers"
)

func main() {
	http.HandleFunc("/api/tcsteps", handlers.TCStepsHandler)
	http.HandleFunc("/api/tcsubmit", handlers.TCSubmitHandler)
	http.HandleFunc("/api/aitnumbers", handlers.AITNumbersHandler)
	http.HandleFunc("/api/expectedoutputs", handlers.ExpectedOutputsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe(":8080", nil)
}
