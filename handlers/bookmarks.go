package handlers

import (
	"github.com/aspersieman/cooper/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Bookmark struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func GetBookmarks(c *gin.Context) {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT * FROM bookmarks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bookmarks []Bookmark
	for rows.Next() {
		var bookmark Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.URL, &bookmark.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		bookmarks = append(bookmarks, bookmark)
	}

	c.JSON(http.StatusOK, bookmarks)
}

func CreateBookmark(c *gin.Context) {
	var bookmark Bookmark
	err := c.BindJSON(&bookmark)
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("INSERT INTO bookmarks (title, url) VALUES (?, ?)", bookmark.Title, bookmark.URL)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, bookmark)
}

func GetBookmark(c *gin.Context) {
	id := c.Param("id")

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	var bookmark Bookmark
	err = dbConn.QueryRow("SELECT * FROM bookmarks WHERE id = ?", id).Scan(&bookmark.ID, &bookmark.Title, &bookmark.URL, &bookmark.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, bookmark)
}

func UpdateBookmark(c *gin.Context) {
	id := c.Param("id")

	var bookmark Bookmark
	err := c.BindJSON(&bookmark)
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("UPDATE bookmarks SET title = ?, url = ? WHERE id = ?", bookmark.Title, bookmark.URL, id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, bookmark)
}

func DeleteBookmark(c *gin.Context) {
	id := c.Param("id")

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("DELETE FROM bookmarks WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusNoContent, nil)
}
