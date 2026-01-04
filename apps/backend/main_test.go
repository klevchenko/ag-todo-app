package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite" // SQLite драйвер
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Ця функція створює тестову базу даних у пам'яті (In-Memory)
func setupTestDB() *gorm.DB {
	// Відкриваємо SQLite в пам'яті (:memory:)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	// Створюємо таблицю
	db.AutoMigrate(&Task{})
	return db
}

// Тест створення задачі
func TestCreateTask(t *testing.T) {
	// 1. Підготовка (Setup)
	db := setupTestDB()
	router := setupRouter(db) // Передаємо тестову базу в наш роутер

	// 2. Створюємо HTTP запит (ніби це Postman)
	newTask := Task{Title: "Test Unit Task"}
	jsonValue, _ := json.Marshal(newTask)
	
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder() // Сюди сервер запише відповідь

	// 3. Виконання (Execute)
	router.ServeHTTP(w, req)

	// 4. Перевірка (Assert)
	// Перевіряємо статус код
	assert.Equal(t, http.StatusCreated, w.Code)

	// Перевіряємо, що повернувся JSON з ID
	var responseTask Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)
	
	assert.Equal(t, "Test Unit Task", responseTask.Title)
	assert.NotZero(t, responseTask.ID) // ID має бути створений
}

// Тест отримання списку
func TestGetTasks(t *testing.T) {
	// 1. Setup
	db := setupTestDB()
	// Запишемо щось в базу напряму перед тестом
	db.Create(&Task{Title: "Existing Task", Done: true})
	
	router := setupRouter(db)

	// 2. Execute
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// 3. Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []Task
	json.Unmarshal(w.Body.Bytes(), &tasks)

	assert.Len(t, tasks, 1) // Має бути 1 задача
	assert.Equal(t, "Existing Task", tasks[0].Title)
	assert.True(t, tasks[0].Done)
}