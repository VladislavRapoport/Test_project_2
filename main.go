package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var messages []string

	// Проверка Load Average
	loadStr := r.URL.Query().Get("loadavg")
	if loadStr != "" {
		if v, err := strconv.ParseFloat(loadStr, 64); err == nil {
			if v > 5.0 {
				messages = append(messages, fmt.Sprintf("Load average high: %.2f", v))
			}
		}
	}

	// Проверка свободного места на диске
	diskStr := r.URL.Query().Get("disk")
	if diskStr != "" {
		if v, err := strconv.Atoi(diskStr); err == nil {
			if v < 10 {
				messages = append(messages, fmt.Sprintf("Free disk space low: %d Gb left", v))
			}
		}
	}

	// Проверка пропускной способности сети
	netStr := r.URL.Query().Get("network")
	if netStr != "" {
		if v, err := strconv.Atoi(netStr); err == nil {
			if v < 20 {
				messages = append(messages, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", v))
			}
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Если сообщений нет, ничего не пишем (тест ожидает пустую строку)
	if len(messages) > 0 {
		_, _ = w.Write([]byte(strings.Join(messages, "\n")))
	}
}

func main() {
	http.HandleFunc("/stats", statsHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}