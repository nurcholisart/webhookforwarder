package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jsonBody map[string]interface{} = map[string]interface{}{}

type Sink struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Sinks []Sink

func main() {
	e := echo.New()

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/", func(c echo.Context) error {
		if err := json.NewDecoder(c.Request().Body).Decode(&jsonBody); err != nil {
			return err
		}

		urlEnv := os.Getenv("URL")
		if urlEnv == "" {
			return c.String(http.StatusOK, "no URL env")
		}

		urls := strings.Split(urlEnv, " ")

		for _, url := range urls {
			go makeHttpRequest(url, jsonBody)
		}

		return c.String(http.StatusOK, "ok")
	})

	port := getEnv("PORT", "8000")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func makeHttpRequest(url string, jsonBody map[string]interface{}) {
	jsonBodyStr, err := json.Marshal(jsonBody)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Forwarding request to %v\n", url)

	http.Post(url, "application/json", bytes.NewBuffer(jsonBodyStr))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
