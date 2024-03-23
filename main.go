package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func selectionSort(arr []int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}
}

func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = key
	}
}

// WebSocket endpoint for handling real-time sorting
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage:", err)
			break
		}

		var request struct {
			Array     []int  `json:"array"`
			Algorithm string `json:"algorithm"`
		}

		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Unmarshal:", err)
			continue
		}

		// Choose the sorting algorithm based on the request
		switch request.Algorithm {
		case "bubble":
			bubbleSort(request.Array)
		case "selection":
			selectionSort(request.Array)
		case "insertion":
			insertionSort(request.Array)
		default:
			log.Println("Invalid algorithm specified")
			continue
		}

		// Send the sorted array back through the WebSocket
		if err := conn.WriteJSON(request.Array); err != nil {
			log.Println("WriteJSON:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/ws", wsEndpoint) // Set up WebSocket route
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
