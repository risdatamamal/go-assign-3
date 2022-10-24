package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	engine "github.com/risdatamamal/go-assign-3/config/gin"
	docs "github.com/risdatamamal/go-assign-3/docs"
	"github.com/risdatamamal/go-assign-3/pkg/domain/message"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func main() {
	// gin engine
	ginEngine := engine.NewGinHttp(engine.Config{
		Port: ":8080",
	})
	// setiap request yang datang ke API ini,
	// dia akan melalui gin.Recovery dan gin.Logger
	// .USE disini, adalah cara untuk memasukkan middleware juga
	ginEngine.GetGin().Use(
		gin.Recovery(),
		gin.Logger())

	startTime := time.Now()
	ginEngine.GetGin().GET("/", func(ctx *gin.Context) {
		// secara default map jika di return dalam
		// response API, dia akan menjadi JSON
		respMap := map[string]any{
			"code":       0,
			"message":    "server up and running",
			"start_time": startTime,
		}

		// golang memiliki json package
		// json package bisa mentranslasikan
		// map menjadi suatu struct
		// nb: struct harus memiliki tag/annotation JSON
		var respStruct message.Response

		// marshal -> mengubah json/struct/map menjadi
		// array of byte atau bisa kita translatekan menjadi string
		// dengan format JSON
		resByte, err := json.Marshal(respMap)
		if err != nil {
			panic(err)
		}
		// unmarshal -> translasikan string/[]byte dengan format JSON
		// menjadi map/struct dengan tag/annotation json
		err = json.Unmarshal(resByte, &respStruct)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, respStruct)
	})

	docs.SwaggerInfo.BasePath = "/v1"
	ginEngine.GetGin().GET("/swagger/*any", ginswagger.
		WrapHandler(swaggerfiles.Handler))

	// for template rendering
	ginEngine.GetGin().LoadHTMLFiles("template/index.html")
	type DataPoint struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	}
	data := []DataPoint{}
	ginEngine.GetGin().GET("/index", func(c *gin.Context) {
		newData := DataPoint{
			Water: rand.Intn(100),
			Wind:  rand.Intn(100),
		}

		latestStatus := "NORMAL"
		if newData.Wind > 8 {
			if newData.Water > 8 {
				latestStatus = "BAHAYA"
			}
		}
		data = append(data, newData)

		file, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile("test.json", file, 0644)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":         "Value website",
			"latest_status": latestStatus,
			"data":          data,
		})
	})
	ginEngine.Serve()
}

func init() {
	godotenv.Load(".env")
}
