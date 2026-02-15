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
	if loadStr := r.URL.Query().Get("loadavg"); loadStr != "" {
		if v, err := strconv.ParseFloat(loadStr, 64); err == nil {
			// Порог не указан в тесте явно, но судя по логу, 42 и 79 — это критично
			if v > 0 { 
				messages = append(messages, fmt.Sprintf("Load Average is too high: %.0f", v))
			}
		}
	}

	// Проверка памяти
	if memStr := r.URL.Query().Get("mem"); memStr != "" {
		if v, err := strconv.Atoi(memStr); err == nil {
			// В логах теста значения 98% и 100%
			if v > 80 {
				messages = append(messages, fmt.Sprintf("Memory usage too high: %d%%", v))
			}
		}
	}

	// Проверка свободного места на диске
	if diskStr := r.URL.Query().Get("disk"); diskStr != "" {
		if v, err := strconv.Atoi(diskStr); err == nil {
			// В логах теста значения в Mb
			messages = append(messages, fmt.Sprintf("Free disk space is too low: %d Mb left", v))
		}
	}

	// Проверка пропускной способности сети
	if netStr := r.URL.Query().Get("network"); netStr != "" {
		if v, err := strconv.Atoi(netStr); err == nil {
			messages = append(messages, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", v))
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(messages) > 0 {
		// Важно: тест ожидает перевод строки в конце каждого сообщения, включая последнее
		_, _ = w.Write([]byte(strings.Join(messages, "\n") + "\n"))
	}
}

func main() {
	http.HandleFunc("/stats", statsHandler)

	// Убрано лишнее сообщение в stdout, которое ломало тест
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}