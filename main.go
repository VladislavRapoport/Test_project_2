package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var messages []string

	if v, err := strconv.Atoi(q.Get("mem")); err == nil && q.Get("mem") != "" {
		messages = append(messages, fmt.Sprintf("Memory usage too high: %d%%", v))
	}
	if v, err := strconv.Atoi(q.Get("loadavg")); err == nil && q.Get("loadavg") != "" {
		messages = append(messages, fmt.Sprintf("Load Average is too high: %d", v))
	}
	if v, err := strconv.Atoi(q.Get("disk")); err == nil && q.Get("disk") != "" {
		messages = append(messages, fmt.Sprintf("Free disk space is too low: %d Mb left", v))
	}
	if v, err := strconv.Atoi(q.Get("network")); err == nil && q.Get("network") != "" {
		messages = append(messages, fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available", v))
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if len(messages) == 0 {
		return
	}

	out := strings.Join(messages, "\n") + "\n"

	// ответ по HTTP
	_, _ = fmt.Fprint(w, out)

	// и обязательно в stdout (именно это сравнивает тест)
	fmt.Print(out)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", statsHandler)

	// Тест, как правило, обращается к srv.msk01.gigacorp.local по 80 порту.
	// Можно слушать просто ":80", этого достаточно.
	_ = http.ListenAndServe(":80", mux)
}