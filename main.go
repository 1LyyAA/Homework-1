package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Status struct {
	Status string `json:"status"`
}

type Message struct {
	Message string `json:"message"`
}

// Handler для GET /
func Gethandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET /] request received")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Пытаемся прочитать из файла конфига
	message, err := os.ReadFile("/app/config/welcome_message")
	if err != nil {
		fmt.Println("[GET /] config file not found, using default")
		fmt.Fprint(w, "Welcome to the custom app")
		return
	}
	fmt.Println("[GET /] returning message from config")
	fmt.Fprint(w, string(message))
}

// Handler для GET /status
func GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET /status] request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	status := Status{Status: "ok"}
	json.NewEncoder(w).Encode(status)
}

// Handler для POST /log
func PostLogHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Println("[POST /log] request received")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		fmt.Printf("[POST /log] JSON decode error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	fmt.Printf("[POST /log] message: %s\n", msg.Message)

	// Создаём директорию если её нет
	logDir := "/app/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("[POST /log] mkdir error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating directory: %v", err)
		return
	}

	// Записываем логи в файл
	logFile := filepath.Join(logDir, "app.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[POST /log] file open error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s\n", msg.Message); err != nil {
		fmt.Printf("[POST /log] write error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error writing to file: %v", err)
		return
	}

	fmt.Println("[POST /log] logged successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"logged"}`)
}

// Handler для GET /logs
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET /logs] request received")
	logFile := "/app/logs/app.log"
	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("[GET /logs] file not found, returning empty")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "")
			return
		}
		fmt.Printf("[GET /logs] read error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading file: %v", err)
		return
	}

	fmt.Printf("[GET /logs] returning %d bytes\n", len(content))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func main() {
	fmt.Println("=== Server starting ===")
	fmt.Println("Server is running on port 8080")
	http.HandleFunc("/", Gethandler)
	http.HandleFunc("/log", PostLogHandler)
	http.HandleFunc("/logs", GetLogsHandler)
	http.HandleFunc("/status", GetStatusHandler)

	fmt.Println("Handlers registered, listening on :8080")
	http.ListenAndServe(":8080", nil)

}
