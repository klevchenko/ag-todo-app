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
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}

// --- SERVER STRUCT ---
// Ми створюємо структуру, яка тримає підключення до бази.
// Тепер наші методи будуть належати цій структурі.
type Server struct {
	DB *gorm.DB
}

// --- HANDLERS (Methods of Server) ---

// Зверни увагу: func (s *Server) ...
// Тепер ми беремо базу не з глобальної змінної, а з s.DB
func (s *Server) getTasks(c *gin.Context) {
	var tasks []Task
	if result := s.DB.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (s *Server) createTask(c *gin.Context) {
	var input Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := Task{Title: input.Title, Done: false}
	if result := s.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (s *Server) updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := s.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	var input Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.DB.Model(&task).Updates(input)
	c.JSON(http.StatusOK, task)
}

func (s *Server) deleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := s.DB.Delete(&Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// --- ROUTER SETUP ---
// Ми винесли налаштування роутера в окрему функцію.
// Це дозволить нам запускати API в тестах без реального запуску сервера.
func setupRouter(db *gorm.DB) *gin.Engine {
	server := &Server{DB: db} // Створюємо сервер з нашою базою
	
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.GET("/tasks", server.getTasks)
		api.POST("/tasks", server.createTask)
		api.PUT("/tasks/:id", server.updateTask)
		api.DELETE("/tasks/:id", server.deleteTask)
	}
	return r
}

// --- MAIN ---
func main() {
	_ = godotenv.Load("/vault/secrets/config")

	// --- ДОДАЛИ ПЕРЕВІРКУ ТУТ ---
	apiKey := os.Getenv("API_KEY")
	fmt.Printf("🔑 Loaded API_KEY: %s\n", apiKey)

	host := "postgres-postgresql"
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable",
		host, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&Task{})

	// Викликаємо нашу функцію налаштування
	r := setupRouter(db)

	fmt.Println("🚀 Server running on :8080")
	r.Run(":8080")
}