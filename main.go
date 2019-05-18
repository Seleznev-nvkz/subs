package main

import (
	"os"

	"github.com/aaaton/golem"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	lemm       *golem.Lemmatizer
	translator *Translator
	db         *DB
	addr       string
)

func init() {
	lemm, _ = golem.New("english")

	yandexAPIKey, ok := os.LookupEnv("YA_API_KEY")
	if !ok {
		panic("Not Found YA_API_KEY")
	}
	translator = newTranslator(yandexAPIKey, "ru")

	addr, ok = os.LookupEnv("ADDR")
	if !ok {
		addr = ":8090"
	}

	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		dbPath = "subtitles.db"
	}
	db = newDB(dbPath)
	db.init()
}

func getRoute() *gin.Engine {
	route := gin.Default()
	v1 := route.Group("/api")
	subtitlesRegister(v1.Group("/subtitles"))
	wordsRegister(v1.Group("/words"))

	route.Use(static.Serve("/", static.LocalFile("./static", true)))
	route.GET("/favicon.ico", func(c *gin.Context) { c.String(200, "ok") })
	route.NoRoute(func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	return route
}

func main() {
	defer db.Close()
	getRoute().Run(addr)
}
