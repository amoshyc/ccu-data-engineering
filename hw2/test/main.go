package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	var inp = flag.String("i", "/dev/stdin", "path/to/input")
	var out = flag.String("o", "/dev/stdout", "path/to/output")
	flag.Parse()

	mysort(*inp, *out)
}

func mysort(inputPath, outputPath string) {
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

	data = parallelSort(data)
	// sort.Slice(data, func(i, j int) bool {
	// 	return data[i] > data[j]
	// })

	fmt.Println(time.Since(startTime))

	for _, line := range data {
		writer.WriteString(line)
		writer.WriteString("\n")
	}
}

func parallelSort(data []string) []string {
	res := make([]string, len(data))
	mergeSort(data, res, 0, len(data), 0)
	return res
}

func mergeSort(data []string, res []string, lb, ub int, dep int) {
	if dep >= 4 {
		sort.Slice(data[lb:ub], func(i, j int) bool {
			return data[lb+i] > data[lb+j]
		})
		for i := lb; i < ub; i++ {
			res[i] = data[i]
		}
		return
	}

	pv := (lb + ub) / 2
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		mergeSort(data, res, lb, pv, dep+1)
		wg.Done()
	}()
	go func() {
		mergeSort(data, res, pv, ub, dep+1)
		wg.Done()
	}()
	wg.Wait()

	nl, nr := pv-lb, ub-pv
	l, r := 0, 0
	for i := lb; i < ub; i++ {
		switch {
		case r == nr:
			res[i] = data[lb+l]
			l++
		case l == nl:
			res[i] = data[pv+r]
			r++
		case data[lb+l] > data[pv+r]:
			res[i] = data[lb+l]
			l++
		default:
			res[i] = data[pv+r]
			r++
		}
	}

	for i := lb; i < ub; i++ {
		data[i] = res[i]
	}
}
