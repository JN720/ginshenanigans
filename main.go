package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type user struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
	Auth     string `json:"auth"`
	Created  string `json:"created"`
}

type newUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
	Auth     string `json:"auth"`
}

func userString(usr *newUser) string {
	return "'" + usr.Email + "', '" + usr.Password + "', '" + usr.Status + "', " + usr.Auth
}

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("environmental variable for POSTGRES_URL not found")
		return
	}
	connStr := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	app := gin.Default()
	app.GET("/", func(c *gin.Context) { get(c, db) })
	app.GET("/:id", func(c *gin.Context) { getUser(c, db) })
	app.POST("/", func(c *gin.Context) { postUser(c, db) })
	app.Run("localhost:8000")
}

func get(c *gin.Context, db *sql.DB) {
	c.Status(http.StatusOK)
}

func getUser(c *gin.Context, db *sql.DB) {
	var usr user
	id, _ := c.Params.Get("id")
	err := db.QueryRow("SELECT * FROM Users WHERE id = "+id+";").Scan(&usr.Id, &usr.Email, &usr.Password, &usr.Status, &usr.Auth, &usr.Created)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, usr)
}

func postUser(c *gin.Context, db *sql.DB) {
	var newUsr newUser

	if err := c.ShouldBindJSON(&newUsr); err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	query := "INSERT INTO Users(email, password, status, auth, created) VALUES(" + userString(&newUsr) + ", NOW()) RETURNING *;"
	fmt.Println(query)
	var usr user

	if err := db.QueryRow(query).Scan(&usr.Id, &usr.Email, &usr.Password, &usr.Status, &usr.Auth, &usr.Created); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}
