package controllers

import (
	"historm_api/models"
	"historm_api/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type DetailBody struct {
	Link string `bson:"link" json:"link"`
}

func scrapAnimeDetail(link string) models.AnimeDetail {
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

	animeType, _ := doc.Find("div.anime_info_body_bg > p.type").First().Find("a").Attr("title")

	var plot string = doc.Find("div.anime_info_body_bg > p.type").Eq(1).Text()
	plot = strings.Trim(strings.Split(plot, ":")[1], " ")

	var genres []string = []string{}
	doc.Find("div.anime_info_body_bg > p.type").Eq(2).Find("a").Each(func(i int, s *goquery.Selection) {
		genre, _ := s.Attr("title")
		genres = append(genres, genre)
	})

	var released string = doc.Find("div.anime_info_body_bg > p.type").Eq(3).Text()
	released = strings.Trim(strings.Split(released, ":")[1], " ")

	status, _ := doc.Find("div.anime_info_body_bg > p.type").Eq(4).Find("a").Attr("title")

	var other string = doc.Find("div.anime_info_body_bg > p.type").Eq(5).Text()
	other = strings.Trim(strings.Split(other, ":")[1], " ")

	epLength, _ := doc.Find("ul#episode_page > li").Last().Find("a").Attr("ep_end")

	var episodes []models.Episode

	length, _ := strconv.Atoi(epLength)

	for i := 1; i <= length; i++ {
		episodes = append(episodes, models.Episode{
			Title: "Episode " + strconv.Itoa(i),
			Link:  utils.BaseUrl + "/" + strings.Split(link, "ry/")[1] + "-episode-" + strconv.Itoa(i),
		})
	}

	return models.AnimeDetail{
		Type:     animeType,
		Plot:     plot,
		Genres:   genres,
		Released: released,
		Status:   status,
		Other:    other,
		Episodes: episodes,
	}
}

func RetrieveAnimeDetail(ctx *gin.Context) {
	var body DetailBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	SendNotification()
	ctx.JSON(http.StatusOK, gin.H{
		"data": scrapAnimeDetail(body.Link),
	})
}
