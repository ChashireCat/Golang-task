package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//	"log"
	"net/http"

	//	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"

	//	"html/template"
	"os"
	//	"os/signal"
	"sync"

	"github.com/gin-gonic/gin"
)

type (
	Auth struct {
		Id       uint64 `json:"Id"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	User struct {
		Id       uint64  `json:"Id"`
		Email    string  `form:"Email" json:"Email"`
		Password string  `form:"Password" json:"Password"`
		balance  float64 `json:"Balance"`
	}
	Response struct {
		Data string `json:"data"`
	}
	Token struct {
		Token string `json:"token"`
	}
)

var (
	cache = make(map[string]*User) // [login]*User
	mut   sync.Mutex
	user  = Auth{
		Id:       1,
		Email:    "email",
		Password: "password",
	}
)

func midle(c *gin.Context) {
	fmt.Println("Ip ", c.ClientIP(), "Serv ", c.Request.URL)
}

func handlerUserCreate(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("handlerUserCreate ERR = ", err)
		return
	}

	var user User

	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Printf("UserCreate err = %s; body = %s\n", err.Error(), string(body))
		return
	}
	val, ok := getUser(user.Email) //Проверка на уникальность *

	if val == nil {
		c.JSON(200, "User register") //Проверка на уникальность *

	} else if ok && val != nil {

		c.JSON(200, "User already register") // User already register
		return
	}

	fmt.Printf("%+v\n", user)
	createUser(user)
}
func createUser(user User) {
	// insert to db

	mut.Lock()
	cache[user.Email] = &user
	mut.Unlock()

}
func getUser(Email string) (*User, bool) {
	mut.Lock()
	val, ok := cache[Email]
	mut.Unlock()

	return val, ok
}
func UserAuthHandler(c *gin.Context) {
	var u Auth
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//сравнение
	if user.Email != u.Email || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
func CreateToken(userid uint64) (string, error) {
	var err error
	//Создание токена
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func main() {
	router := gin.New()

	router.Use(midle)
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	router.POST("/register", handlerUserCreate)

	router.POST("/UserAuth", UserAuthHandler)

	router.GET("/authorize", func(c *gin.Context) {
		c.HTML(http.StatusOK, "authorize.html", gin.H{
			"title": "authorize",
		})
	})
	router.GET("/account", func(c *gin.Context) {
		c.HTML(http.StatusOK, "account.html", gin.H{
			"title": "account",
		})
	})

	router.GET("/homepage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homepage.html", gin.H{
			"title": "Home Page",
		})
	})
	router.Run(":8080")
}
