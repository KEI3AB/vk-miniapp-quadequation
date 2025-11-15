package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
)

type quadRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

type quadResponse struct {
	Roots []float64 `json:"roots"`
	Msg   string    `json:"message,omitempty"`
}

func quadCalc(a, b, c float64) ([]float64, error) {
	if a == 0 {
		if b == 0 {
			return nil, errors.New("invalid input (a = 0, b = 0)")
		}
		return []float64{-c / b}, nil
	}

	d := b*b - 4*a*c
	if d < 0 {
		return nil, errors.New("отрицательный дискриминант")
	} else if d == 0 {
		return []float64{-b / (2 * a)}, nil
	} else {
		return []float64{(-b + math.Sqrt(d)) / (2 * a),
			(-b - math.Sqrt(d)) / (2 * a)}, nil
	}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var req quadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	roots, err := quadCalc(req.A, req.B, req.C)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quadResponse{Roots: roots})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/solve", solveHandler)

	fmt.Println("Server is running on the port: ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server error: ", err)
	}
}
