package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://hypeauditor.com/top-instagram-all-russia/")
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to fetch the page"))
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to parse the HTML"))
		return
	}

	file, err := os.Create("task3/users.csv")
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed create csv fail"))
		return
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"рейтинг", "ник", "имя", "категории", "последователи", "страна", "слушатели (Auth.)", "слушатели (Avg.)"})

	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to write CSV headers"))

	}
	doc.Find(".row__top").Each(func(i int, s *goquery.Selection) {

		rating := s.Find(".rank span[data-v-65566291]").Text()
		name := s.Find(".contributor__name").Text()
		username := s.Find(".contributor__title").Text()
		category := s.Find(".category[data-v-65566291]").Text()
		followers := s.Find(".subscribers[data-v-65566291]").Text()
		caunter := s.Find(".audience[data-v-3c683b7e]").Text()
		audeenter := s.Find(".authentic[data-v-3c683b7e]").Text()
		audeen_Avg := s.Find(".engagement[data-v-3c683b7e]").Text()

		row := []string{rating, name, username, category, followers, caunter, audeenter, audeen_Avg}
		err = writer.Write(row)

		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to write date to CSV"))
		}
	})

	fmt.Println("Data has to task3/users.csv")
}
