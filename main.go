package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	// initialize Echo
	e := echo.New()

	// Set middleware
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Route requests
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/:id", func(c echo.Context) error {
		if err := json.NewDecoder(c.Request().Body).Decode(&jsonBody); err != nil {
			return err
		}

		jsonFile, err := os.Open("sinks.json")
		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		var sinks Sinks

		byteJsonFile, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJsonFile, &sinks)

		sinkId := c.Param("id")
		for _, sink := range sinks {
			if sink.ID == sinkId {
				go makeHttpRequest(sink.URL, jsonBody)
			}
		}

		return c.String(http.StatusOK, "ok")
	})

	e.Logger.Fatal(e.Start(":8000"))
}

func makeHttpRequest(url string, jsonBody map[string]interface{}) {
	jsonBodyStr, err := json.Marshal(jsonBody)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("Forwarding request to %v\n", url)

	http.Post(url, "application/json", bytes.NewBuffer(jsonBodyStr))
	return
}
