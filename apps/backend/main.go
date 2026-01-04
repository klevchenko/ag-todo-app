package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Шлях, куди Vault покладе файл.
// Це стандартний шлях: /vault/secrets/<ім'я-секрету-в-конфігу-кубернетису>
const secretPath = "/vault/secrets/config"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Читаємо файл, який створив Vault Agent
		content, err := os.ReadFile(secretPath)
		if err != nil {
			fmt.Fprintf(w, "Error reading secret: %v\n", err)
			return
		}

		// 2. Виводимо вміст (для навчання). 
		// У реальному житті ми б використали це для підключення до БД.
		fmt.Fprintf(w, "Hello from Go!\nSecret received from Vault:\n%s", string(content))
	})

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}