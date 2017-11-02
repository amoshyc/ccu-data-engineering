package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Item ...
type Item struct {
	s    string
	freq int
}

func main() {
	var inp = flag.String("i", "/dev/stdin", "path/to/input")
	var out = flag.String("o", "/dev/stdout", "path/to/output")
	flag.Parse()

	myCount(*inp, *out)
}

func myCount(inputPath, outputPath string) {
	fin, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	fout, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	scanner := bufio.NewScanner(fin)
	data := make([]string, 0)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	startTime := time.Now()

	cnt := make(map[string]int)
	for _, line := range data {
		line = strings.Replace(line, "\t", " ", -1)
		tokens := strings.Split(line, " ")
		for _, token := range tokens {
			if _, exist := cnt[token]; exist {
				cnt[token]++
			} else {
				cnt[token] = 1
			}
		}
	}

	res := make([]Item, 0)
	for k, v := range cnt {
		res = append(res, Item{k, v})
	}
	sort.Slice(res, func(i, j int) bool {
		freq1, freq2 := res[i].freq, res[j].freq
		if freq1 == freq2 {
			return res[i].s < res[j].s
		}
		return freq1 > freq2
	})

	fmt.Println(time.Since(startTime))

	for _, item := range res {
		output := fmt.Sprintf("%d %s", item.freq, item.s)
		writer.WriteString(output)
		writer.WriteString("\n")
	}
}
