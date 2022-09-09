package controllers

import (
	"historm_api/models"
	"historm_api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Body struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
	Page  int    `bson:"page" json:"page"`
}

func scrapAnimeList(animeType string, value string, page int) []models.Anime {
	var url string = ""
	var baseUrl string = utils.BaseUrl
	if animeType == "dub" {
		url = baseUrl + "/search.html?keyword=dub&page=" + strconv.Itoa(page)
	} else if animeType == "genre" {
		url = baseUrl + "/genre/" + strings.Join(strings.Split(strings.ToLower(value), " "), "-") + "?page=" + strconv.Itoa(page)
	} else if animeType == "type" {
		url = baseUrl + "/" + strings.Join(strings.Split(strings.ToLower(value), " "), "-") + ".html?page=" + strconv.Itoa(page)
	} else if animeType == "subcategory" {
		url = baseUrl + "/sub-category/" + strings.Join(strings.Split(strings.ToLower(value), " "), "-") + "?page=" + strconv.Itoa(page)
	} else {
		url = ""
	}

	var c *colly.Collector = colly.NewCollector()
	var animeList []models.Anime

	c.OnHTML("div.main_body > div.last_episodes > ul.items > li", func(e *colly.HTMLElement) {
		var anime models.Anime = models.Anime{
			Title: e.ChildAttr("p.name, a", "title"),
			Image: e.ChildAttr("img", "src"),
			Link:  baseUrl + e.ChildAttr("p.name, a", "href"),
		}
		animeList = append(animeList, anime)
	})

	c.Visit(url)
	return animeList
}

func RetrieveAnimeList(ctx *gin.Context) {
	var body Body
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	SendNotification()

	ctx.JSON(http.StatusOK, gin.H{
		"data": scrapAnimeList(body.Type, body.Value, body.Page),
	})
}
