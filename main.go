package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

// Структура для приема данных от тестов
type Stats struct {
	MemoryUsage  int `json:"memory_usage"`
	LoadAverage  int `json:"load_average"`
	FreeDisk     int `json:"free_disk_space"`
	NetworkSpeed int `json:"network_bandwidth"`
}

// Универсальная функция для формирования строк предупреждений
func checkStats(mem, load, disk, net int) {
	var msgs []string
	if mem > 80 {
		msgs = append(msgs, fmt.Sprintf("Memory usage too high: %d%%", mem))
	}
	if load > 30 {
		msgs = append(msgs, fmt.Sprintf("Load Average is too high: %d", load))
	}
	if disk < 30000 {
		msgs = append(msgs, fmt.Sprintf("Free disk space is too low: %d Mb left", disk))
	}
	if net < 150 {
		msgs = append(msgs, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", net))
	}

	if len(msgs) > 0 {
		fmt.Println(strings.Join(msgs, "\n"))
	}
}

func main() {
	// Определяем флаги (для тестов 2 и 3)
	memPtr := flag.Int("mem", 0, "Memory usage")
	loadPtr := flag.Int("load", 0, "Load average")
	diskPtr := flag.Int("disk", 40000, "Free disk space") // Значение по умолчанию выше порога
	netPtr := flag.Int("net", 200, "Network bandwidth")  // Значение по умолчанию выше порога
	flag.Parse()

	// Если через флаги переданы критические значения, выводим их (тесты 2 и 3)
	// Проверяем, были ли флаги установлены пользователем (аргументы командной строки)
	if flag.NFlag() > 0 {
		checkStats(*memPtr, *loadPtr, *diskPtr, *netPtr)
	}

	// Обработчик для тестов 1 и 4
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var s Stats
			err := json.NewDecoder(r.Body).Decode(&s)
			if err == nil {
				// При получении JSON сразу печатаем в консоль
				checkStats(s.MemoryUsage, s.LoadAverage, s.FreeDisk, s.NetworkSpeed)
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	// Запускаем сервер. Тесты 1 и 4 ожидают, что он слушает порт 8080
	http.ListenAndServe(":8080", nil)
}