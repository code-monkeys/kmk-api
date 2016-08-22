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
    "encoding/json"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

var (
    repeat int
    db *sql.DB
)

type Todo struct {
    id      int    // `json:"id"`
    user_id string // `json:"user_id"`
    content string // `json:"content"`
}

type Todos struct {
    list[] *Todo
}

// func newTodos() *Todos {

// }
// func (g *Todos) addTodo(user_id int, content string) {

// }

func repeatHandler(c *gin.Context) {
    var buffer bytes.Buffer
    for i := 0; i < repeat; i++ {
        buffer.WriteString("Hello from Go!\n")
    }
    c.String(http.StatusOK, buffer.String())
}

func execDb(c *gin.Context, query string) {
    if _, err := db.Exec(query); err != nil {
        c.String(http.StatusInternalServerError,
            fmt.Sprintf("Error creating database table: %q", err))
        return
    }
}

func getDb(c *gin.Context, query string) {
    rows, err := db.Query("SELECT tick FROM ticks")
    if err != nil {
        c.String(http.StatusInternalServerError,
            fmt.Sprintf("Error reading ticks: %q", err))
        return
    }

    defer rows.Close()
    for rows.Next() {
        
    }
}

func initDb(c *gin.Context) {
    execDb(c, "CREATE TABLE IF NOT EXISTS ticks (tick timestamp)")
    execDb(c, "INSERT INTO ticks VALUES (now())")

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
    defer db.Close()

    router := gin.New()
    router.Use(gin.Logger())
    router.LoadHTMLGlob("templates/*.tmpl.html")
    router.Static("/static", "static")

    router.GET("/", func(c *gin.Context) {
        // c.HTML(http.StatusOK, "index.tmpl.html", nil)
        data := []string{"apple", "peach", "pear"}
        list, _ := json.Marshal(data)
        bool, _ := json.Marshal(true)
        int1, _ := json.Marshal(123)
        flot, _ := json.Marshal(123.456)
        str , _ := json.Marshal("23dd")

        c.JSON(200, gin.H{
            "name":    "kmk-api",
            "version": "v0.1.0--",
            "list": list,
            "bool": bool,
            "int" : int1,
            "float" : flot,
            "str" : str,
        })
    })

    // router.GET("/info", func(c *gin.Context) {
    //  c.HTML(http.StatusOK, "index.tmpl.html", nil)
    // })

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    router.GET("/todos/:user", func(c *gin.Context) {
        user := c.Param("user")
        c.String(http.StatusOK, "Hello %s", user)
    })

    router.GET("/todos/:user/:id", func(c *gin.Context) {
        user := c.Param("user")
        id   := c.Param("id")
        c.JSON(200, gin.H{
            "user": user,
            "id":   id,
        })
    })

    router.Run(":" + port)
}
