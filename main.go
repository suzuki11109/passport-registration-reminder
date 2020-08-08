package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	formURL = "https://docs.google.com/forms/d/e/1FAIpQLSeGvkG2Eoe8d0FVNxxH5GL9HbXFUgysVTtK9XLs_4rYKqPXxQ/viewform"
	word    = "มีผู้ลงทะเบียนเต็มแล้ว"
)

func main() {
	collector := colly.NewCollector()

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	userID := os.Getenv("USER_ID")
	if userID == "" {
		log.Fatal(errors.New("user id is required"))
	}

	fmt.Println("starting...")

	t := time.NewTicker(time.Duration(3 * time.Minute))
	defer t.Stop()

	go func() {
		for t := range t.C {
			fmt.Println("run at", t)
			if !checkRegistrationForm(collector) {
				messageMe(bot, userID, "Registration is ready! Let's do it right away")
				messageMe(bot, userID, formURL)
			} else {
				fmt.Println("Still not ready for registration")
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	fmt.Println("done")
}

func messageMe(bot *linebot.Client, userID, text string) {
	message := linebot.NewTextMessage(text)

	if _, err := bot.PushMessage(userID, message).Do(); err != nil {
		fmt.Println(err)
	}
}

func checkRegistrationForm(collector *colly.Collector) bool {
	var found bool
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	collector.OnHTML(fmt.Sprintf("[data-value='%s']", word), func(e *colly.HTMLElement) {
		found = true
	})

	collector.Visit(formURL)
	return found
}
