package main

// 2.6GB, 1m1s

import (
	"bufio"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type Item struct {
	Term string
	Pos []([2]int)
}

func NewItem(term string) *Item {
	res := new(Item)
	res.Term = term
	res.Pos = make([][2]int, 0)
	return res
}

func (item *Item) Add(pos [2]int) {
	item.Pos = append(item.Pos, pos);
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ngrams(line_id int, lines []string, terms *[]string) map[string](*Item) {
	res := make(map[string]*Item)
	for _, term := range *terms {
		res[term] = NewItem(term)
	}

	for row_id, line := range lines {
		indices := make([]int, 0)
		for i, w := 0, 0; i < len(line); i += w {
			_, width := utf8.DecodeRuneInString(line[i:])
			indices = append(indices, i)
			w = width
		}

		N := len(indices)
		s, t := 0, min(7, N)
		indices = append(indices, len(line))

		for s < N-1 {
			query := line[indices[s]:indices[t]]
			if _, exist := res[query]; exist {
				res[query].Add([2]int{line_id + row_id, s})
				s, t = t, min(t+7, N)
			} else {
				if t-s == 2 {
					s, t = s+1, min(s+8, N)
				} else {
					t = t - 1
				}
			}
		}
	}

	return res
}

func read_file(path string) []string {
	content, _ := ioutil.ReadFile(path)
	data := strings.Split(string(content), "\n")
	return data
}

func input() string {
	reader := bufio.NewReader(os.Stdin)
	inp, _ := reader.ReadString('\n')
	return inp
}

func st_ngrams(text_path, term_path, output_path string) []Item {
	terms := read_file(term_path)
	text := read_file(text_path)
	res := ngrams(0, text, &terms)

	merged := make(map[string]*Item)
	for _, term := range terms {
		merged[term] = NewItem(term)
	}

	for k, item := range res {
		for _, val := range item.Pos {
			merged[k].Add(val)
		}
	}

	result := make([]Item, 0)
	for _, item := range merged {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		len_i := len(result[i].Pos)
		len_j := len(result[j].Pos)
		if len_i == len_j {
			return result[i].Term < result[j].Term
		}
		return len(result[i].Pos) > len(result[j].Pos)
	})

	json_res, _ := json.Marshal(result)
	ioutil.WriteFile(output_path, json_res, 0644)

	return result
}


func main() {
	// text := []string{"中a文b也c", "中ab也", "ababab"}
	// term := []string{"中a", "b也", "ab"}
	// fmt.Println(ngrams(text, term))

	// text_path := "../../data/10/wiki_00"
	// term_path := "../../data/terms.txt"
	// out_path := "./out.json"

	// text_path := "../../data/pu/doc.txt"
	// term_path := "../../data/pu/term.txt"

	// mt_ngrams(text_path, term_path)

	term_path := os.Args[1]
	text_path := os.Args[2]
	out_path := os.Args[3]
	st_ngrams(text_path, term_path, out_path)
}
