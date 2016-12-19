package lm

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	comment    = regexp.MustCompile("\\\\.*\\\\")
	newNgram   = regexp.MustCompile("\\\\\\d+")
	contentRow = regexp.MustCompile("(-)?\\d+(\\.\\d+)?\\s.*?((-)?\\d+(\\.\\d+))?")
	backoffRow = regexp.MustCompile(".*?(-)?\\d+(\\.\\d+)")
	whiteSpace = regexp.MustCompile("\\s+")
)

func Build(fileName string) LanguageModel {
	var lm LanguageModel = make(map[string]word)
	file, err := os.Open(fileName)
	check(err, "File could not be opened.", 0)
	defer file.Close()
	scan := bufio.NewScanner(file)
	row := 0
	ngram := 0
	for scan.Scan() {
		text := scan.Text()
		row += 1
		switch {
		case newNgram.MatchString(text):
			ngram += 1
		case contentRow.MatchString(text):
			list := whiteSpace.Split(text, -1)
			if len(list) < ngram+1 {
				panic("File is not correctly formated")
			}
			lm.addWord(list, row, ngram)
		}
	}
	NGRAM = ngram
	return lm
}

func (lm *LanguageModel) addWord(row []string, r, ngram int) {
	prob, err := strconv.ParseFloat(row[0], 64)
	check(err, "File is not correctly formated", r)
	var backoff float64 = 0
	bo, err := strconv.ParseFloat(row[len(row)-1], 64)
	if err == nil {
		backoff = bo
	}
	w := word{prob, backoff}
	(*lm)[strings.TrimSpace(strings.Join(row[1:ngram+1], " "))] = w
}

func check(e error, s string, row int) {
	if e != nil {
		log.Println(e)
		log.Fatal("Error on row ", row, ". ", s)
	}
}
