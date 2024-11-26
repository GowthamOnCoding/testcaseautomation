package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Step struct {
	StepName       string `json:"step_name"`
	TCStep         string `json:"tcstep"`
	Value          string `json:"value"`
	ExpectedOutput string `json:"expected_output"` // Added expected output field
}

type TestCase struct {
	TestCaseName string `json:"test_case_name"`
	AITNumber    string `json:"aitnumber"` // Added AIT number field
	Steps        []Step `json:"steps"`
}

func TCSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var testCase TestCase
	err := json.NewDecoder(r.Body).Decode(&testCase)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process the test case here (e.g., save to database)
	fmt.Printf("Received Test Case: %+v\n", testCase)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test case submitted successfully!"))
}
