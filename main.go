package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jumadimuhammad/backendgo/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func app(e *echo.Echo, store model.UserStore) {

	e.POST("/register", func(c echo.Context) error {
		name := c.FormValue("name")
		address := c.FormValue("address")
		telp, _ := strconv.Atoi(c.FormValue("telp"))
		email := c.FormValue("email")
		password := c.FormValue("password")
		role := "3"
		token := "secret"

		checkemail := store.FindEmail(email)

		if checkemail != nil {
			return echo.ErrUnauthorized
		}

		if password == "" {
			return echo.ErrUnauthorized
		}

		hashpwd, _ := model.Hash(password)

		user, _ := model.CreateUser(name, address, telp, email, hashpwd, role, token)

		store.Save(user)

		return c.JSON(http.StatusOK, user)
	})

	e.POST("/login", func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		if password == "" || email == "" {
			return echo.ErrUnauthorized
		}

		user := store.Login(email)

		err := model.CheckPasswordHash(password, user.Password)

		if err != true {
			return echo.ErrUnauthorized
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["name"] = user.Name
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		t, _ := token.SignedString([]byte("secret"))

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	})

	e.GET("/", func(c echo.Context) error {
		users := store.All()

		return c.JSON(http.StatusOK, users)
	})

	r := e.Group("/users")
	r.Use(middleware.JWT([]byte("secret")))

	r.GET("/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		user := store.Find(id)

		return c.JSON(http.StatusOK, user)
	})

	r.GET("/:role/role", func(c echo.Context) error {
		role, _ := strconv.Atoi(c.Param("role"))

		users := store.FindRole(role)

		return c.JSON(http.StatusOK, users)
	})

	r.PUT("/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		user := store.Find(id)
		user.Name = c.FormValue("name")
		user.Address = c.FormValue("address")
		user.Telp, _ = strconv.Atoi(c.FormValue("telp"))
		user.Email = c.FormValue("email")
		password := c.FormValue("password")

		user.Password, _ = model.Hash(password)

		store.Update(user)

		return c.JSON(http.StatusOK, user)
	})

	r.DELETE("/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		user := store.Find(id)

		store.Delete(user)

		return c.JSON(http.StatusOK, user)
	})

}

func main() {

	godotenv.Load()
	var store model.UserStore
	store = model.NewUserMySQL()

	e := echo.New()
	e.Use(middleware.CORS())

	app(e, store)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
