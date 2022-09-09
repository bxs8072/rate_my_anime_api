package controllers

import (
	"historm_api/models"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type VideoBody struct {
	Link string `bson:"link" json:"link"`
}

func scrapeVideo(link string) models.Video {

	resp, err := http.Get(link)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	iframe, _ := doc.Find("div.play-video > iframe").Attr("src")

	download, _ := doc.Find("div.anime_video_body_cate > div.favorites_book > ul > li.dowloads > a ").Attr("href")

	return models.Video{
		Iframe:   "https:" + iframe,
		Download: download,
	}
}

func RetrieveVideo(ctx *gin.Context) {
	var body VideoBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": scrapeVideo(body.Link),
	})
}
