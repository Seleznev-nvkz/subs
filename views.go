package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func subtitlesRegister(router *gin.RouterGroup) {
	router.DELETE("/:slug", subtitlesDelete)
	router.GET("/", subtitlesList)
	router.GET("/:slug/translate", subtitlesTranslate)
}

func wordsRegister(router *gin.RouterGroup) {
	router.DELETE("/:slug", wordDelete)
	router.POST("/upload", uploadFile)
}

func wordDelete(c *gin.Context) {
	exception := Exceptions{Word: c.Param("slug")}
	err := db.Save(&exception)
	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
	} else {
		c.String(http.StatusNoContent, "")
	}
}

func subtitlesTranslate(c *gin.Context) {
	var subtitles Subtitles

	slug, err := strconv.Atoi(c.Param("slug"))
	if err == nil {
		err = db.One("ID", slug, &subtitles)
		if err != nil {
			c.String(http.StatusBadRequest, "%s", err)
		} else {
			for i := range subtitles.Words {
				subtitles.Words[i].translate()
			}
			db.Save(&subtitles)
			c.JSON(http.StatusOK, subtitles)
		}
	} else {
		c.String(http.StatusBadRequest, "%s", err)
	}
}

// defaultMultipartMemory = 32 << 20 // 32 MB
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.String(http.StatusBadRequest, "%s", err)
	} else {
		src, err := file.Open()
		defer src.Close()

		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err)
		}

		wordsArray, err := handleFile(src)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err)
		}

		words := make([]Word, len(wordsArray))
		for i := range wordsArray {
			words[i] = Word{Word: wordsArray[i]}
		}
		subtitles := Subtitles{Words: words, Name: file.Filename}
		db.Save(&subtitles)

		c.JSON(http.StatusOK, subtitles)
	}
}

func subtitlesList(c *gin.Context) {
	var subtitles []Subtitles

	db.AllByIndex("ID", &subtitles)
	for i := range subtitles {
		subtitles[i].refresh()
	}
	c.JSON(http.StatusOK, subtitles)
}

func subtitlesDelete(c *gin.Context) {
	var subtitles Subtitles

	slug, err := strconv.Atoi(c.Param("slug"))
	if err == nil {
		err = db.One("ID", slug, &subtitles)
		if err != nil {
			c.String(http.StatusBadRequest, "%s", err)
		} else {
			db.DeleteStruct(&subtitles)
			c.String(http.StatusNoContent, "")
		}
	} else {
		c.String(http.StatusBadRequest, "%s", err)
	}
}
