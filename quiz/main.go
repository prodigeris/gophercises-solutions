package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	f := flag.String("filename", "problems.csv", "filename of problems file")
	flag.Parse()

	data := readCsv("quiz/" + *f)
	c := 0
	for _, row := range data {
		q := row[0]
		a := row[1]
		fmt.Println(q)
		var i string
		_, err := fmt.Scanln(&i)
		if err != nil {
			i = ""
		}
		if i == a {
			c++
		}
	}

	fmt.Println("You got", c, "out of", len(data))
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
