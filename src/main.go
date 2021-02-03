package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"restfuluser/src/data"
	"restfuluser/src/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupRoutes(e *echo.Echo) {
	e.GET("/user/:uid", handlers.GetUserByID)
	e.GET("/user", handlers.GetUsers)
	e.POST("/user", handlers.CreateUser)
	e.PUT("/user/:uid", handlers.UpdateUserByID)
	e.DELETE("/user/:uid", handlers.DeleteUserByID)
}

func loadUsers() (*data.Users, error) {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	jsonFile.Close()
	users := new(data.Users)
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uint64]data.User)
	for _, u := range users.Users {
		usersMap[u.UID] = u
	}
	users.UsersMap = usersMap
	return users, nil
}

func main() {
	server := echo.New()
	users, err := loadUsers()
	if err != nil {
		server.Logger.Fatal(err)
	}

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			uc := &data.UserContext{ctx, users}
			return next(uc)
		}
	})

	server.Use(handlers.RejectNonJSONRequests)

	server.Use(
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			if len(reqBody) != 0 {
				fmt.Printf(" with payload:\n  %v", string(reqBody))
			}
			fmt.Print("\n")
		}),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${method} @ ${uri} => ${status}",
		}))

	setupRoutes(server)
	server.Logger.Fatal(server.Start(":1323"))
}
