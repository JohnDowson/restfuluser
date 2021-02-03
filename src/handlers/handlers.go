package handlers

import (
	"fmt"
	"net/http"
	"restfuluser/src/data"
	"strconv"

	"github.com/labstack/echo/v4"
)

// RejectNonJSONRequests rejects all GET requests with Content-Type not matching application/json
func RejectNonJSONRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Request().Method != http.MethodGet {
			isJSON := ctx.Request().Header["Content-Type"][0] == "application/json"
			if isJSON {
				return next(ctx)
			}
			return ctx.String(http.StatusBadRequest, "400 Bad Request: Bad Content-Type")
		}
		return next(ctx)

	}
}

type errorResponce struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func notFoundError(uid uint64) errorResponce {
	return errorResponce{
		http.StatusNotFound,
		"Not Found",
		fmt.Sprintf("No user with UID %d exists", uid)}
}
func badRequestError(detail string) errorResponce {
	return errorResponce{
		http.StatusBadRequest,
		"Bad Request",
		detail}
}
func internalServerError(detail string) errorResponce {
	return errorResponce{
		http.StatusInternalServerError,
		"Internal Server Error",
		detail}
}

// GetUserByID shut
func GetUserByID(c echo.Context) error {
	uid, err := uidFromParam(c)
	if err != nil {
		return err
	}
	uc := c.(*data.UserContext)
	user := uc.Users.Get(uid)
	if user != nil {
		return c.JSON(http.StatusOK, user)
	}
	return c.JSON(http.StatusNotFound, notFoundError(uid))
}

// GetUsers shut
func GetUsers(c echo.Context) error {
	uc := c.(*data.UserContext)
	return c.JSON(http.StatusOK, uc.Users.Users)
}

// CreateUser shut
func CreateUser(c echo.Context) error {
	uc := c.(*data.UserContext)
	iu := new(data.IncompleteUser)
	err := c.Bind(iu)
	if err != nil {
		return err
	}
	u, err := uc.Users.Insert(*iu)
	if err != nil {
		fmt.Printf("Following error has occured when creating user: \n  %v\n", err)
		return c.JSON(http.StatusInternalServerError, internalServerError(""))
	}
	return c.JSON(http.StatusOK, u)
}

// UpdateUserByID shut
func UpdateUserByID(c echo.Context) error {
	uid, err := uidFromParam(c)
	if err != nil {
		return err
	}
	uc := c.(*data.UserContext)
	iu := new(data.IncompleteUser)
	err = c.Bind(iu)
	if err != nil {
		return err
	}
	user, err := uc.Users.Update(uid, *iu)
	if err != nil {
		fmt.Printf("Following error has occured when updating user: \n  %v\n", err)
		return c.JSON(http.StatusInternalServerError, internalServerError(""))
	}
	if user != nil {
		return c.JSON(http.StatusOK, user)
	}
	return c.JSON(http.StatusNotFound, notFoundError(uid))
}

// DeleteUserByID shut
func DeleteUserByID(c echo.Context) error {
	uid, err := uidFromParam(c)
	if err != nil {
		return err
	}
	uc := c.(*data.UserContext)
	err = uc.Users.Delete(uid)
	if err != nil {
		fmt.Printf("Following error has occured when deleting user: \n  %v\n", err)
		return c.JSON(http.StatusInternalServerError, internalServerError(""))
	}
	return c.NoContent(http.StatusOK)
}

func uidFromParam(c echo.Context) (uint64, error) {
	maybeUID := c.Param("uid")
	uid, err := strconv.ParseUint(maybeUID, 10, 64)
	if err != nil {
		errDetail := fmt.Sprintf("'%s' is not a valid UID", maybeUID)
		return 0, c.JSON(http.StatusBadRequest, badRequestError(errDetail))
	}
	return uid, nil
}
