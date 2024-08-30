package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`(\p{L}[\p{L}|-]*)|(-{2,})`)

type item struct {
	Str   string
	Count int
}

type (
	hashStatistic map[string]int
	wordStatistic []item
)

func Top10(input string) []string {
	return parse(input).sort().top10()
}

func parse(input string) wordStatistic {
	h := make(hashStatistic)

	for _, s := range reg.FindAllString(input, -1) {
		h[strings.ToLower(s)]++
	}
	return h.toWordStatictic()
}

func (h hashStatistic) toWordStatictic() wordStatistic {
	s := make(wordStatistic, 0, len(h))
	for k, v := range h {
		s = append(s, item{Str: k, Count: v})
	}
	return s
}

func (wordStat wordStatistic) sort() wordStatistic {
	sort.Slice(wordStat, func(i, j int) bool {
		return wordStat[i].Count > wordStat[j].Count ||
			(wordStat[i].Count == wordStat[j].Count && wordStat[i].Str < wordStat[j].Str)
	})

	return wordStat
}

func (wordStat wordStatistic) top10() []string {
	results := make([]string, 0, 10)

	for i := 0; i < 10 && i < len(wordStat); i++ {
		results = append(results, wordStat[i].Str)
	}
	return results
}
