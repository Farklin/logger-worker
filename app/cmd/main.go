package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

type Logger struct {
	Mutex sync.Mutex
	Data  map[string]int
}

func main() {

	logger := newLogger()
	log := &Logger{
		Data: make(map[string]int),
	}

	go logging(logger)
	go loggingSecond(log)

	http.HandleFunc("/log", func(w http.ResponseWriter, request *http.Request) {
		file, err := os.ReadFile("activity.log")
		if err != nil {
			fmt.Fprintf(w, "Ошибка чтения файла "+err.Error())
		}
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, string(file))
	})

	http.HandleFunc("/log-second", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(log)
		if err != nil {
			fmt.Errorf("Ошибка формата ")
		}
		fmt.Fprintf(w, string(data))
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logger.Println("Ошибка запуска сервера")
	}

}

func newLogger() *log.Logger {
	// Открываем файл для записи логов
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}

	logger := log.New(file, "", log.LstdFlags)
	return logger
}

func loggingSecond(log *Logger) {
	for {
		windowTitle := robotgo.GetTitle()
		log.Mutex.Lock()
		log.Data[windowTitle]++
		log.Mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func logging(logger *log.Logger) {

	for {
		// Получаем активное окно
		windowTitle := robotgo.GetTitle()
		logger.Printf("Активное окно: %s", windowTitle)

		// Получаем список процессов
		//procs, _ := process.Processes()
		//for _, p := range procs {
		//	name, _ := p.Name()
		//	if len(name) > 0 {
		//		logger.Printf("Процесс: %s (PID: %d)", name, p.Pid)
		//	}
		//}

		// Задержка перед следующим логированием
		time.Sleep(5 * time.Second)
	}
}
