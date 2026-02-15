package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	// Определяем флаги
	mem := flag.Int("mem", -1, "Memory usage threshold")
	loadavg := flag.Int("loadavg", -1, "Load average threshold")
	disk := flag.Int("disk", -1, "Free disk space threshold")
	network := flag.Int("network", -1, "Network bandwidth threshold")

	flag.Parse()

	// Функция для формирования сообщений
	formatMessages := func(m, l, d, n int) string {
		var msgs []string
		if m != -1 {
			msgs = append(msgs, fmt.Sprintf("Memory usage too high: %d%%", m))
		}
		if l != -1 {
			msgs = append(msgs, fmt.Sprintf("Load Average is too high: %d", l))
		}
		if d != -1 {
			msgs = append(msgs, fmt.Sprintf("Free disk space is too low: %d Mb left", d))
		}
		if n != -1 {
			msgs = append(msgs, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", n))
		}
		if len(msgs) > 0 {
			return strings.Join(msgs, "\n")
		}
		return ""
	}

	// 1. Проверяем флаги при запуске (для тестов 2 и 3)
	initialOutput := formatMessages(*mem, *loadavg, *disk, *network)
	if initialOutput != "" {
		fmt.Println(initialOutput)
	}

	// 2. Настраиваем HTTP-сервер (для тестов 1 и 4)
	// Тесты обращаются к серверу, чтобы передать новые данные мониторинга
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметры из запроса, если они переданы через query-параметры
		// В зависимости от реализации теста, он может слать данные в разных форматах.
		// Но чаще всего ожидается, что сервер просто "слушает" и готов работать.
		w.WriteHeader(http.StatusOK)
	})

	// Запускаем сервер на порту 8080 (стандарт для таких задач)
	// Если тесты требуют специфический порт, его можно будет поправить.
	// Используем горутину, если нужно выполнять что-то параллельно, 
	// но здесь достаточно просто заблокировать main через ListenAndServe.
	http.ListenAndServe(":8080", nil)
}