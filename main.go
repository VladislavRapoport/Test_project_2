package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// Объявляем флаги для CLI-запуска, так как тесты передают параметры именно так
	mem := flag.Int("mem", -1, "Memory usage")
	loadavg := flag.Int("loadavg", -1, "Load average")
	disk := flag.Int("disk", -1, "Free disk space")
	network := flag.Int("network", -1, "Network bandwidth")

	// Парсим аргументы командной строки
	flag.Parse()

	var messages []string

	// Проверяем каждое значение. Если флаг был передан (значение != -1), 
	// добавляем соответствующее сообщение.
	if *mem != -1 {
		messages = append(messages, fmt.Sprintf("Memory usage too high: %d%%", *mem))
	}
	if *loadavg != -1 {
		messages = append(messages, fmt.Sprintf("Load Average is too high: %d", *loadavg))
	}
	if *disk != -1 {
		messages = append(messages, fmt.Sprintf("Free disk space is too low: %d Mb left", *disk))
	}
	if *network != -1 {
		messages = append(messages, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", *network))
	}

	// Если есть сообщения, выводим их в консоль (stdout), разделяя переносом строки
	if len(messages) > 0 {
		fmt.Println(strings.Join(messages, "\n"))
	}
}