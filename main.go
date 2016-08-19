package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"bytes"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
	// "github.com/russross/blackfriday"
	_ "github.com/lib/pq"
)

var (
	repeat int
	db *sql.DB
)

func repeatHandler(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!\n")
	}
	c.String(http.StatusOK, buffer.String())
}

func dbFunc(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
		  c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error scanning ticks: %q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}


func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err = strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}


	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		// c.HTML(http.StatusOK, "index.tmpl.html", nil)
		c.JSON(200, gin.H{
			"name":    "kmk-api",
			"version": "v0.1.0",
		})
	})

	// router.GET("/info", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/todos/:user", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.GET("/todos/:user/:id", func(c *gin.Context) {
		name := c.Param("name")
		id   := c.Param("id")
		c.JSON(200, gin.H{
			"name": name,
			"id":   id,
		})
	})

	// router.GET("/md", func(c *gin.Context) {
	// 	c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	// })

	// router.GET("/yo", repeatHandler)
	// router.GET("/db", dbFunc)

	router.Run(":" + port)
}
