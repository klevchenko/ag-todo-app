package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- MODELS ---

type Task struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title" binding:"required"` // binding:"required" валідує вхідний JSON
	Done  bool   `json:"done"`
}

// --- GLOBAL DB ---
var db *gorm.DB

// --- HANDLERS ---

// GET /tasks
func getTasks(c *gin.Context) {
	var tasks []Task
	// Знайти всі задачі
	result := db.Find(&tasks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// POST /tasks
func createTask(c *gin.Context) {
	var input Task
	// Валідація JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Створення в БД
	task := Task{Title: input.Title, Done: false}
	result := db.Create(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// PUT /tasks/:id
func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	// 1. Шукаємо задачу
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 2. Читаємо нові дані
	var input Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Оновлюємо
	db.Model(&task).Updates(input)
	c.JSON(http.StatusOK, task)
}

// DELETE /tasks/:id
func deleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// --- MAIN ---

func main() {
	// 1. Config & Secret Loading
	_ = godotenv.Load("/vault/secrets/config") // Ігноруємо помилку, якщо файлу немає (для локального тесту)

	host := "postgres-postgresql"
	// Якщо ми запускаємо локально (не в k8s), можна переозначити хост через ENV, наприклад localhost
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable",
		host,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	// 2. Database Connection
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// У проді краще падати, якщо бази немає
		panic("Failed to connect to database: " + err.Error())
	}

	// 3. Migration
	db.AutoMigrate(&Task{})

	// 4. Router Setup (Gin)
	r := gin.Default()
	
	// Групуємо API версіювання - це Best Practice
	api := r.Group("/api/v1")
	{
		api.GET("/tasks", getTasks)
		api.POST("/tasks", createTask)
		api.PUT("/tasks/:id", updateTask)
		api.DELETE("/tasks/:id", deleteTask)
	}

	fmt.Println("🚀 Server running on :8080")
	r.Run(":8080")
}