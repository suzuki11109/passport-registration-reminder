package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

const word = "มีผู้ลงทะเบียนเต็มแล้ว"

func main() {
	var found bool
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(fmt.Sprintf("[data-value='%s']", word), func(e *colly.HTMLElement) {
		found = true
	})

	c.Visit("https://docs.google.com/forms/d/e/1FAIpQLSeGvkG2Eoe8d0FVNxxH5GL9HbXFUgysVTtK9XLs_4rYKqPXxQ/viewform")

	if !found {
		fmt.Println("Registration is ready! Let's do it right away")
	} else {
		fmt.Println("Still not ready for registration")
	}
}
