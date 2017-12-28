package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s := newSolver()
	for {
		fmt.Print("> ")

		var inp string
		n, err := fmt.Scanln(&inp)
		if err != nil || n == 0 {
			break
		}

		cs, ok := s.ProcessInput(inp)
		fmt.Println(ok, cs)
	}
}

type candidate string
type candidates []candidate

type solver struct {
	db map[string]candidates
}

func newSolver() *solver {
	s := new(solver)
	s.db = make(map[string]candidates)

	fdb, err := os.Open("data/liu.csv")
	if err != nil {
		panic(err)
	}
	defer fdb.Close()

	csv := csv.NewReader(fdb)
	csv.Comma = '\t'
	records, err := csv.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		word, code := record[0], record[1]
		s.db[code] = append(s.db[code], candidate(word))
	}

	return s
}

func (s *solver) QueryCode(query string) (candidates, bool) {
	val, exist := s.db[query]
	return val, exist
}

func (s *solver) ProcessInput(input string) (candidates, bool) {
	idx := strings.IndexAny(input, "0123456789")
	if idx == -1 {
		return s.QueryCode(input)
	}

	code := input[:idx]
	cs, ok := s.QueryCode(code)
	num, err := strconv.Atoi(input[idx:])
	if err != nil || !ok || num < 1 || num > len(cs) {
		return cs, false
	}
	return cs[num-1 : num], true
}
