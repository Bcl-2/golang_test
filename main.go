package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type News struct {
	ID         int    `db:"Id" json:"Id"`
	Title      string `db:"Title" json:"Title"`
	Content    string `db:"Content" json:"Content"`
	Categories []int  `json:"Categories"`
}

var db *sqlx.DB

func init() {
	// Загрузка значений из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Инициализация подключения к базе данных
	initDB()

	router := gin.Default()

	// Обработчики маршрутов
	router.POST("/edit/:Id", editNews)
	router.GET("/list", getNewsList)

	// Запуск сервера
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server is running on port %s", port)
	log.Fatal(router.Run(":" + port))
}

func initDB() {
	var err error

	// Параметры подключения к базе данных
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Установка соединения с базой данных
	db, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Установка connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}

// Функция изменения новости
func editNews(c *gin.Context) {
	// Получение значения параметра Id из URL
	id, err := strconv.Atoi(c.Param("Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}

	// Получение данных из тела запроса
	var news News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Вызов хранимой процедуры для изменения новости в таблице News
	if _, err := db.Exec("CALL content.UpdateNews(?, ?, ?)", id, news.Title, news.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	// Парсинг идентификаторов из массива Categories
	var stringNumbers []string
	for _, num := range news.Categories {
		stringNumbers = append(stringNumbers, strconv.Itoa(num))
	}
	result := strings.Join(stringNumbers, ",")

	// Изменение идентификаторов в таблице NewsCategories
	if _, err := db.Exec("CALL content.UpdateCategory(?, ?)", id, result); err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getNewsList(c *gin.Context) {
	// Вызов хранимой процедуры для получения списка новостей
	newsResp, err := db.Queryx("CALL content.GetNews()")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news list"})
		return
	}
	defer newsResp.Close()

	// Сбор данных новостей
	var newsList []News
	for newsResp.Next() {
		var news News
		if err := newsResp.StructScan(&news); err != nil {
			log.Println(err)
			continue
		}

		categories, err := db.Queryx("CALL content.GetCategories(?)", news.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		defer categories.Close()

		// Сбор данных категорий
		var id int
		for categories.Next() {
			if err := categories.Scan(&id); err != nil {
				log.Println(err)
				continue
			}
			news.Categories = append(news.Categories, id)
		}
		newsList = append(newsList, news)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "News": newsList})
}
