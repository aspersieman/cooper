package main

import (
	"github.com/aspersieman/cooper/handlers"
	"log"
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // Import the SQLite driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // Import the file source driver
	_ "github.com/mattn/go-sqlite3"
)

//go:embed static
var staticFS embed.FS

func main() {
	doMigrate()
	// Create a new Gin router
	r := gin.Default()

	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("static/html/index.htm", http.FS(staticFS))
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("web/img/favicon.ico", http.FS(staticFS))
	})
	cssFS, err := fs.Sub(staticFS, "static/css")
	if err != nil {
		panic(err)
	}
	jsFS, err := fs.Sub(staticFS, "static/js")
	if err != nil {
		panic(err)
	}
	imgFS, err := fs.Sub(staticFS, "static/img")
	if err != nil {
		panic(err)
	}

	r.StaticFS("/css", http.FS(cssFS))
	r.StaticFS("/js", http.FS(jsFS))
	r.StaticFS("/img", http.FS(imgFS))

	// Define API endpoints
	r.GET("/bookmarks", handlers.GetBookmarks)
	r.POST("/bookmarks", handlers.CreateBookmark)
	r.GET("/bookmarks/:id", handlers.GetBookmark)

	r.PUT("/bookmarks/:id", handlers.UpdateBookmark)
	r.DELETE("/bookmarks/:id", handlers.DeleteBookmark)
	
	port := ":8122"
	log.Printf("Server running on port %s", port)
	// Start the server
	r.Run(port)
}

func doMigrate() {
	// Open a database connection (or create if it doesn't exist)
	db, err := sql.Open("sqlite3", "./my_database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize golang-migrate
	m, err := migrate.New(
		"file://./db/migrations",
		"sqlite3://./storage/db/cooper.db",
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Apply all available migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Database migrations applied successfully!")
}
