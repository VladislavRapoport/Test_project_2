package main

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    fmt.Println("🖥️  Мониторинг сервера запущен...")
    
    for {
        // Проверка CPU (имитация 20-80%)
        cpuUsage := getCPUUsage()
        fmt.Printf("📊 CPU: %d%%\n", cpuUsage)
        
        // Проверка памяти (имитация 40-90%)
        memUsage := getMemoryUsage()
        fmt.Printf("💾 RAM: %d%% (%.1f GB)\n", memUsage, float64(memUsage)/100*8)
        
        // Проверка дискового пространства (60-95%)
        diskUsage := getDiskUsage()
        fmt.Printf("💿 Диск: %d%%\n", diskUsage)
        
        fmt.Printf("⏰ %s\n", time.Now().Format("15:04:05"))
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println()
        
        time.Sleep(3 * time.Second)
    }
}

func getCPUUsage() int {
    return 20 + rand.Intn(60)
}

func getMemoryUsage() int {
    return 40 + rand.Intn(50)
}

func getDiskUsage() int {
    return 60 + rand.Intn(36)
}