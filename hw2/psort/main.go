package main

import (
	"fmt"
	"sort"
)

func parallelSort(data []string) []string {
	in := make(chan []string, 1)
	out := make(chan []string, 1)
	go mergesort(in, out, 0)
	in <- data
	return <-out
}

func mergesort(in chan []string, out chan []string, dep int) {
	data := <-in

	if dep >= 2 {
		go func() {
			sort.Slice(data, func(i, j int) bool {
				return data[i] > data[j]
			})
		}()
		out <- data
		return
	}

	N := len(data)
	in1 := make(chan []string, 1)
	in2 := make(chan []string, 1)
	res1 := make(chan []string, 1)
	res2 := make(chan []string, 1)
	mergesort(in1, res1, dep+1)
	mergesort(in2, res2, dep+1)
	in1 <- data[:N/2]
	in2 <- data[N/2:]

	l, r := <-res1, <-res2
	i, j := 0, 0
	res := make([]string, N)
	for ix := range res {
		switch {
		case i == len(l):
			res[ix] = r[j]
			j++
		case j == len(r):
			res[ix] = l[i]
			i++
		case l[i] > r[j]:
			res[ix] = l[i]
			i++
		default:
			res[ix] = r[j]
			j++
		}
	}

	out <- res
}

func main() {
	data := []string{
		"136",
		"751",
		"951",
		"136",
		"751",
		"951",
	}

	data = parallelSort(data)
	fmt.Println(data)
}
