package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

func main() {
	var inp = flag.String("i", "/dev/stdin", "path/to/input")
	var out = flag.String("o", "/dev/stdout", "path/to/output")
	var mxn = flag.Int("n", 5, "max n")
	flag.Parse()

	for i := 2; i <= *mxn; i++ {
		fmt.Print(i, ": ")
		words := ngramsAndCount(*inp, i)
		writeWords(*out, i, words)
		fmt.Println("Done")
	}
}

func ngramsAndCount(inputPath string, n int) map[string]int {
	fin, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	isStopWords := map[rune]bool{
		rune('，'): true,
		rune('。'): true,
		rune('！'): true,
		rune('？'): true,
		rune('；'): true,
		rune('：'): true,
		rune('「'): true,
		rune('」'): true,
		rune('（'): true,
		rune('）'): true,
		rune('／'): true,
		rune('、'): true,
		rune('《'): true,
		rune('》'): true,
		rune('〈'): true,
		rune('〉'): true,
	}

	words := make(map[string]int)

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		line := scanner.Text()
		validWords := make([]rune, 0)

		for _, r := range line { // iterate all runes
			if r >= 256 && !isStopWords[r] { // discard ascii & stopwords
				validWords = append(validWords, r)
			}
		}

		for i := 0; i+n < len(validWords); i++ {
			word := string(validWords[i : i+n])
			words[word]++
		}
	}

	return words
}

func writeWords(outputPath string, n int, words map[string]int) {
	outputPath = fmt.Sprintf("%s.%d", outputPath, n)
	fout, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	type item struct {
		k string
		v int
	}
	items := make([]item, 0)
	for k, v := range words {
		items = append(items, item{k, v})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].v == items[j].v {
			return items[i].k < items[j].k
		}
		return items[i].v > items[j].v
	})

	for _, x := range items {
		out := fmt.Sprintf("%s;%09d\n", x.k, x.v)
		writer.WriteString(out)
	}
}
