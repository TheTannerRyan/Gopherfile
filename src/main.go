package main

import (
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type response struct {
	ID        string `json:"id"`
	Success   bool   `json:"success"`
	Payload   string `json:"payload"`
	UserAgent string `json:"userAgent"`
}

func main() {
	var err error
	e := echo.New()

	// index
	e.GET("/", func(c echo.Context) error {
		u := uuid.Must(uuid.NewV4(), err)
		hash := fmt.Sprintf("%s", u)

		return c.JSON(http.StatusOK, &response{
			ID:        hash,
			Success:   true,
			Payload:   time.Now().UTC().Format(time.RFC3339),
			UserAgent: c.Request().UserAgent(),
		})
	})

	e.Start(":3000")
}
