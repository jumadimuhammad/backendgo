package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jumadimuhammad/backendgo/model"
	"github.com/labstack/echo"
)

func app(e *echo.Echo, store model.UserStore) {

	// curl http://localhost:8080
	e.GET("/", func(c echo.Context) error {
		// Process
		user := "Welcome......."

		// Response
		return c.JSON(http.StatusOK, user)
	})

	// curl http://localhost:8080/users
	e.GET("/users", func(c echo.Context) error {
		// Process
		users := store.All()

		// Response
		return c.JSON(http.StatusOK, users)
	})

	// curl http://localhost:8080/users/1
	e.GET("/users/:id", func(c echo.Context) error {
		// Given
		id, _ := strconv.Atoi(c.Param("id"))

		// Process
		user := store.Find(id)

		// Response
		return c.JSON(http.StatusOK, user)
	})

	// curl http://localhost:8080/users/3/role
	e.GET("/users/:role/role", func(c echo.Context) error {
		// Given
		role, _ := strconv.Atoi(c.Param("role"))

		// Process
		users := store.FindRole(role)

		// Response
		return c.JSON(http.StatusOK, users)
	})

	e.POST("/users", func(c echo.Context) error {
		// Given
		name := c.FormValue("name")
		address := c.FormValue("address")
		telp, _ := strconv.Atoi(c.FormValue("telp"))
		email := c.FormValue("email")
		password := c.FormValue("password")
		role := "3"
		token := "secret"

		//Hashing password
		hash, err := model.Hash(password)
		if err != nil {
			log.Fatal(err)
		}

		hashpwd := string(hash)

		// Create instabce
		user, _ := model.CreateUser(name, address, telp, email, hashpwd, role, token)

		// Persist
		store.Save(user)

		// Response
		return c.JSON(http.StatusOK, user)
	})

	e.PUT("/users/:id", func(c echo.Context) error {
		// Given
		id, _ := strconv.Atoi(c.Param("id"))

		// Process
		user := store.Find(id)
		user.Name = c.FormValue("name")
		user.Address = c.FormValue("address")
		user.Telp, _ = strconv.Atoi(c.FormValue("telp"))
		user.Email = c.FormValue("email")
		password := c.FormValue("password")

		hash, err := model.Hash(password)
		if err != nil {
			log.Fatal(err)
		}

		hashpwd := string(hash)

		user.Password = hashpwd

		// Persists
		store.Update(user)

		// Response
		return c.JSON(http.StatusOK, user)
	})

	e.DELETE("/users/:id", func(c echo.Context) error {
		// Given
		id, _ := strconv.Atoi(c.Param("id"))

		// Process
		user := store.Find(id)

		// Remove
		store.Delete(user)

		// Response
		return c.JSON(http.StatusOK, user)
	})

	e.POST("/login", func(c echo.Context) error {
		// Given
		email := c.FormValue("email")
		password := c.FormValue("password")

		user := store.Login(email)

		//Hashing password
		match := model.CheckPasswordHash(password, user.Password)

		// Response
		return c.JSON(http.StatusOK, match)
	})

}

func main() {

	godotenv.Load()
	var store model.UserStore
	store = model.NewUserMySQL()

	e := echo.New()
	app(e, store)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
