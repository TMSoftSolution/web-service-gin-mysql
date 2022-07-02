package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	port     = "3306"
	user     = "root"
	password = ""
	database = "test"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var dbClient sql.DB

// getAlbums responds with the list of all albums as JSON
func getAlbums(c *gin.Context) {

	results, err := dbClient.Query("SELECT * FROM articles")

	if err != nil {
		panic(err.Error())
	}

	var albums = []Album{}
	for results.Next() {
		var album Album

		results.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		albums = append(albums, album)
	}

	defer results.Close()

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the table
	sql := fmt.Sprintf("INSERT INTO articles(title, artist, price) VALUES ('%s', '%s', '%f')", newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	result, err := dbClient.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	result.Close()

	c.IndentedJSON(http.StatusCreated, result)
}

func getAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err.Error())
	}

	var album Album
	sql := fmt.Sprintf("SELECT * FROM articles where id=%d", id)
	err = dbClient.QueryRow(sql).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		panic(err.Error())
	}

	c.IndentedJSON(http.StatusOK, album)

}

func main() {
	dbConnect()
	routers()
}

func routers() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func dbConnect() {
	sqlInfo := fmt.Sprintf("%s:@tcp(%s:%s)/%s", user, host, port, database)
	db, err := sql.Open("mysql", sqlInfo)
	if err != nil {
		panic(err.Error())
	}

	dbClient = *db

	fmt.Println("DB Connected Succesfully.")

	defer db.Close()
}
