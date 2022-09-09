package controllers

import (
	"context"
	"fmt"
	"historm_api/configs"
	"log"
	"regexp"
	"strings"

	"firebase.google.com/go/v4/messaging"
)

var storedTitle string = ""
var storedTitleDub string = ""
var storedTitleChinese string = ""

func SendNotification() {
	firebaseApp, err := configs.FirebaseApp()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fcmClient, err := firebaseApp.Messaging(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}

	animeList := scrapLatestAnimeList("sub", 1)
	lowerTitle := strings.ToLower(animeList[0].Title)
	title := strings.Join(strings.Split(lowerTitle, " "), "")
	title = strings.Trim(re.ReplaceAllString(title, ""), " ")

	animeListDub := scrapLatestAnimeList("dub", 1)
	lowerTitleDub := strings.ToLower(animeListDub[0].Title)
	titleDub := strings.Join(strings.Split(lowerTitleDub, " "), "")
	titleDub = strings.Trim(re.ReplaceAllString(titleDub, ""), " ")

	animeListChinese := scrapLatestAnimeList("chinese", 1)
	lowerTitleChinese := strings.ToLower(animeListChinese[0].Title)
	titleChinese := strings.Join(strings.Split(lowerTitleChinese, " "), "")
	titleChinese = strings.Trim(re.ReplaceAllString(titleChinese, ""), " ")

	if strings.Compare(storedTitle, title) != 0 {
		storedTitle = title
		response, _ := fcmClient.Send(context.TODO(), &messaging.Message{
			Notification: &messaging.Notification{
				Title: "A nice notification title",
				Body:  "A nice notification body",
			},
			Condition: "'" + title + "' in topics" + " || 'all' in topics",
		})
		fmt.Println(response)
	} else if strings.Compare(storedTitleDub, titleDub) != 0 {
		storedTitleDub = titleDub

		response, _ := fcmClient.Send(context.TODO(), &messaging.Message{
			Notification: &messaging.Notification{
				Title: "A nice notification title",
				Body:  "A nice notification body",
			},
			Condition: "'" + titleDub + "' in topics" + " || 'all' in topics",
		})
		fmt.Println(response)

	} else if strings.Compare(storedTitleChinese, titleChinese) != 0 {
		storedTitleChinese = titleChinese

		response, _ := fcmClient.Send(context.TODO(), &messaging.Message{
			Notification: &messaging.Notification{
				Title: "A nice notification title",
				Body:  "A nice notification body",
			},
			Topic: "'" + titleChinese + "' in topics" + " || 'all' in topics",
		})
		fmt.Println(response)
	}
}
