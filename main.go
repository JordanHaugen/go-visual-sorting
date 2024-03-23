package main

import (
	"encoding/json"
	"net/http"
)

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var arr []int
	if err := json.NewDecoder(r.Body).Decode(&arr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bubbleSort(arr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(arr)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/sort", sortHandler)
	http.ListenAndServe(":8080", nil)
}
