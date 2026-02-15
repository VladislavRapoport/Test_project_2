package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func getLoadAverage() string {
	// Использование os.ReadFile вместо ioutil.ReadFile
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "unknown"
	}
	fields := strings.Fields(string(data))
	if len(fields) > 0 {
		return fields[0]
	}
	return "unknown"
}

func getFreeDiskSpace(path string) uint64 {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0
	}
	// Свободные блоки * размер блока = байты
	return stat.Bavail * uint64(stat.Bsize) / 1024 / 1024 // в Мб
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Примерные расчеты для демонстрации формата
	memUsage := (float64(m.Alloc) / float64(m.Sys)) * 100
	diskFree := getFreeDiskSpace("/")
	loadAvg := getLoadAverage()

	result := fmt.Sprintf("Memory usage too high: %.0f%%\n", memUsage)
	result += fmt.Sprintf("Free disk space is too low: %d Mb left\n", diskFree)
	result += fmt.Sprintf("Load Average is too high: %s\n", loadAvg)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
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
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")