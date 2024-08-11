package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	// Conectar a la base de datos MySQL
	dsn := "samuel:21051994@tcp(127.0.0.1:3307)/temperature"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()

	// Definir el endpoint POST para insertar datos
	r.POST("/register", func(c *gin.Context) {
		var input struct {
			Temperature float64 `json:"temperature" binding:"required"`
			Humidity    float64 `json:"humidity" binding:"required"`
		}
		log.Print("ENTRO")
		now := time.Now().Format("2006-01-02 15:04:05")

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insertar datos en la base de datos
		_, err := db.Exec("INSERT INTO registre (date, temperature, humidity, place) VALUES (?, ?, ?, ?)", now, input.Temperature, input.Humidity, "Fuera")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
		
	})


	r.GET("/getAllRegister", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM registre")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var records []struct {			
			Date        string  `json:"date"`
			Temperature float64 `json:"temperature"`
			Humidity    float64 `json:"humidity"`
			Place       string  `json:"place"`
		}

		for rows.Next() {
			var record struct {				
				Date        string  `json:"date"`
				Temperature float64 `json:"temperature"`
				Humidity    float64 `json:"humidity"`
				Place       string  `json:"place"`
			}

			if err := rows.Scan(&record.Date, &record.Temperature, &record.Humidity, &record.Place); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			records = append(records, record)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, records)
	})

	r.Run(":8080") // Ejecutar el servidor en el puerto 8080
}
