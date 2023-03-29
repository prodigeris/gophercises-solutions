package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	f, l := parseFlags()
	data := readCsv(*f)
	c := 0
	fmt.Println("You have", *l, "seconds to answer", len(data), "questions. Press enter to start.")
	_, _ = fmt.Scanln()

	ch := time.NewTimer(time.Duration(*l) * time.Second).C
	c = askQuestions(ch, data, c)

	fmt.Println("\nYou got", c, "out of", len(data))
}

func askQuestions(ch <-chan time.Time, data [][]string, c int) int {
	for _, row := range data {
		q := row[0]
		a := row[1]
		fmt.Println(q + " = ?")

		answerCh := make(chan string)
		go func() {
			var i string
			_, err := fmt.Scanln(&i)
			if err != nil {
				i = ""
			}
			answerCh <- i
		}()
		select {
		case <-ch:
			return c
		case i := <-answerCh:
			if i == a {
				c++
			}
		}
	}
	return c
}

func parseFlags() (*string, *int) {
	f := flag.String("filename", "problems.csv", "filename of problems file")
	l := flag.Int("limit", 30, "default time limit per question")
	flag.Parse()
	return f, l
}

func readCsv(fn string) [][]string {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	return data
}
