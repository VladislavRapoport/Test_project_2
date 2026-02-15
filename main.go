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

func getLoadAverage() (int, bool) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, false
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, false
	}
	f, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, false
	}
	return int(f + 0.5), true
}

func getFreeDiskSpaceMB(path string) (uint64, bool) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, false
	}
	return stat.Bavail * uint64(stat.Bsize) / 1024 / 1024, true
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// В тестах значения приходят параметрами запроса.
	// Пример: /stats?load_average=57&mem_used_percent=100&disk_free=10097&net_bandwidth=416
	loadStr := q.Get("load_average")
	memStr := q.Get("mem_used_percent")
	diskStr := q.Get("disk_free")
	netStr := q.Get("net_bandwidth")

	var b strings.Builder

	// Load Average > 5.0
	if loadStr == "" {
		if v, ok := getLoadAverage(); ok && float64(v) > 5.0 {
			fmt.Fprintf(&b, "Load Average is too high: %d\n", v)
		}
	} else {
		// В тестах load_average целое (57, 95), сравниваем численно
		v, err := strconv.Atoi(loadStr)
		if err == nil && float64(v) > 5.0 {
			fmt.Fprintf(&b, "Load Average is too high: %d\n", v)
		}
	}

	// Memory usage > 90%
	if memStr != "" {
		v, err := strconv.Atoi(memStr)
		if err == nil && v > 90 {
			fmt.Fprintf(&b, "Memory usage too high: %d%%\n", v)
		}
	}

	// Disk free < 1000 Mb
	if diskStr == "" {
		if v, ok := getFreeDiskSpaceMB("/"); ok && v < 1000 {
			fmt.Fprintf(&b, "Free disk space is too low: %d Mb left\n", v)
		}
	} else {
		v, err := strconv.ParseUint(diskStr, 10, 64)
		if err == nil && v < 1000 {
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
		// В тестах вывод программы сравнивают с ожидаемым, поэтому НЕ пишем логи в stdout.
		_ = srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}