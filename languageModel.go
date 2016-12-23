package lm

import (
	"fmt"
	"strings"
)

var NGRAM = 0

type LanguageModel map[string]word

type word struct {
	prob, backoff float64
}

type sentence []string

func ngrams(s sentence, ngram int) [][]string {
	var ngrams [][]string
	if len(s) < ngram {
		return [][]string{s}
	}
	for i := NGRAM; i <= len(s); i++ {
		ngrams = append(ngrams, s[i-ngram:i])
	}
	return ngrams
}

func (l *LanguageModel) Score(s sentence) float64 {
	ngram := ngrams(s, NGRAM)
	sum := 0.0
	for _, s := range ngram {
		ss := l.probNgram(s)
		sum += ss
	}
	return sum
}

func (l *LanguageModel) ScoreUNK(s sentence) float64 {
	return l.Score(l.replaceUNK(s))
}

func (l *LanguageModel) replaceUNK(s sentence) sentence {
	for i := 0; i < len(s); i++ {
		if _, ok := (*l)[s[i]]; !ok {
			s[i] = "<UNK>"
		}
	}
	return s
}

func (l *LanguageModel) ScoreString(s string) float64 {
	ss := append(append([]string{"<s>"}, strings.Split(s, " ")...), "</s>")
	return l.ScoreUNK(ss)
}

func (l *LanguageModel) probNgram(ngram []string) float64 {
	s := strings.TrimSpace(strings.Join(ngram, " "))
	w, ok := (*l)[s]
	if ok {
		fmt.Println("prob: ", w.prob, " for string ", ngram)
		return w.prob
	}
	firstPart := ngram[:len(ngram)-1]
	lastPart := ngram[1:]
	return l.probNgram(lastPart) + l.backoffNgram(firstPart)
}

func (l *LanguageModel) backoffNgram(ngram []string) float64 {
	s := strings.Join(ngram, " ")
	w, ok := (*l)[s]
	if ok {
		return w.backoff
	}
	return 0.0
}
