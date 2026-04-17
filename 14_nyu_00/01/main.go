package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("👷 Trabajador %d iniciando procesamiento paralelo...\n", id)
	time.Sleep(time.Second)
	fmt.Printf("✅ Trabajador %d finalizado\n", id)
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg) // Esto lanza el proceso en paralelo (Goroutine)
	}
	wg.Wait()
	fmt.Println("🚀 Todos los procesos en paralelo han terminado.")
}
