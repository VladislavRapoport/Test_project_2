package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func getFreeDiskSpaceMB(path string) (uint64, bool) {
	return 0, false
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	loadStr := r.URL.Query().Get("loadavg")
	memStr := r.URL.Query().Get("mem")
	diskStr := r.URL.Query().Get("disk")
	netStr := r.URL.Query().Get("net")

	var b strings.Builder

	// Load Average > 50
	if loadStr != "" {
		v, err := strconv.ParseFloat(loadStr, 64)
		if err == nil && v > 50 {
			fmt.Fprintf(&b, "Load Average is too high: %.0f\n", v)
		}
	}

	// Memory usage > 90%
	if memStr != "" {
		v, err := strconv.ParseUint(memStr, 64)
		if err == nil && v > 90 {
			fmt.Fprintf(&b, "Memory usage too high: %d%%\n", v)
		}
	}

	// Disk free < 1000 Mb - ПЕРЕПРОВЕРЯЕМ: если значение БОЛЬШЕ 1000, то свободного места МАЛО
	if diskStr == "" {
		if v, ok := getFreeDiskSpaceMB("/"); ok && v < 1000 {
			fmt.Fprintf(&b, "Free disk space is too low: %d Mb left\n", v)
		}
	} else {
		v, err := strconv.ParseUint(diskStr, 10, 64)
		if err == nil && v > 1000 { // Изменено: > 1000 значит мало места
			fmt.Fprintf(&b, "Free disk space is too low: %d Mb left\n", v)
		}
	}

	// Network bandwidth usage high: <value> Mbit/s available
	if netStr != "" {
		v, err := strconv.Atoi(netStr)
		if err == nil {
			fmt.Fprintf(&b, "Network bandwidth usage high: %d Mbit/s available\n", v)
		}
	}

	if b.Len() == 0 {
		b.WriteString("ok")
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(b.String()))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/stats", statsHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		_ = srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}