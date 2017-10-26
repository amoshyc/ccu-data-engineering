package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func main() {
	N := int(4e7)
	arr := make([]string, N)
	for ix := range arr {
		arr[ix] = strconv.Itoa(N - ix)
	}
	data := make([]string, N)

	copy(data, arr)
	fmt.Print("psort1:")
	st := time.Now()
	ch := make(chan []string, 1)
	mergesort2(data, ch, 0)
	_ = <-ch
	fmt.Println(time.Since(st))

	copy(data, arr)
	fmt.Print("psort2:")
	st = time.Now()
	mergesort1(data, ch, 0)
	_ = <-ch
	fmt.Println(time.Since(st))

	copy(data, arr)
	fmt.Print("sort:")
	st = time.Now()
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})
	fmt.Println(time.Since(st))
}

func mergesort1(data []string, out chan []string, dep int) {
	if dep >= 3 {
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		out <- data
		return
	}

	N := len(data)
	res1 := make(chan []string, 1)
	res2 := make(chan []string, 1)
	go mergesort1(data[:N/2], res1, dep+1)
	go mergesort1(data[N/2:], res2, dep+1)

	res := make([]string, N)
	l, r := <-res1, <-res2
	i, j := 0, 0
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

func mergesort2(data []string, out chan []string, dep int) {
	if dep >= 3 {
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		out <- data
		return
	}

	N := len(data)
	res1 := make(chan []string, 1)
	res2 := make(chan []string, 1)
	go mergesort2(data[:N/2], res1, dep+1)
	go mergesort2(data[N/2:], res2, dep+1)

	l, r := <-res1, <-res2
	i, j := 0, 0
	for ix := range data {
		switch {
		case i == len(l):
			data[ix] = r[j]
			j++
		case j == len(r):
			data[ix] = l[i]
			i++
		case l[i] > r[j]:
			data[ix] = l[i]
			i++
		default:
			data[ix] = r[j]
			j++
		}
	}

	out <- data
}
