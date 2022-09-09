package controllers

import (
	"historm_api/models"
	"historm_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type LatestBody struct {
	Type string `bson:"type" json:"type"`
	Page int    `bson:"page" json:"page"`
}

func scrapLatestAnimeList(animeType string, page int) []models.LatestAnime {
	var t int = 1

	if animeType == "sub" {
		t = 1
	} else if animeType == "dub" {
		t = 2
	} else {
		t = 3
	}

	var url string = "https://ajax.apimovie.xyz/ajax/page-recent-release.html?page=" + strconv.Itoa(page) + "&type=" + strconv.Itoa(t)

	var c *colly.Collector = colly.NewCollector()
	var animeList []models.LatestAnime

	c.OnHTML("div.last_episodes > ul.items > li", func(e *colly.HTMLElement) {
		var anime models.LatestAnime = models.LatestAnime{
			Title: e.ChildAttr("p.name, a", "title"),
			Image: e.ChildAttr("img", "src"),
			Link:  utils.BaseUrl + e.ChildAttr("p.name, a", "href"),
		}
		animeList = append(animeList, anime)
	})

	c.Visit(url)
	return animeList
}

func RetrieveLatestAnimeList(ctx *gin.Context) {
	var body LatestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	SendNotification()
	ctx.JSON(http.StatusOK, gin.H{
		"data": scrapLatestAnimeList(body.Type, body.Page),
	})
}
