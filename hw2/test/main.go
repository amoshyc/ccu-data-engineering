package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	var inp = flag.String("i", "/dev/stdin", "path/to/input")
	var out = flag.String("o", "/dev/stdout", "path/to/output")
	var alg = flag.String("a", "s", "algorithm")
	flag.Parse()

	var startTime time.Time

	if *alg == "p" {
		fmt.Print("psort:")
		startTime = time.Now()
		psort(*inp, *out)
		fmt.Println(time.Since(startTime))
	} else {
		fmt.Print("sort :")
		startTime = time.Now()
		mysort(*inp, *out)
		fmt.Println(time.Since(startTime))
	}
}

func psort(inputPath, outputPath string) {
	fin, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	fout, err := os.Create(outputPath + ".psort")
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

	inps := make([]chan []string, 4)
	outs := make([]chan []string, 4)
	m1 := make(chan []string, 1)
	m2 := make(chan []string, 1)
	res := make(chan []string, 1)

	for ix := range inps {
		inps[ix] = make(chan []string, 1)
		outs[ix] = make(chan []string, 1)
		go func(in, out chan []string) {
			chunk := <-in
			sort.Slice(chunk, func(i, j int) bool {
				return chunk[i] > chunk[j]
			})
			out <- chunk
		}(inps[ix], outs[ix])
	}

	go merge(outs[0], outs[1], m1)
	go merge(outs[2], outs[3], m2)
	go merge(m1, m2, res)

	N := len(data)
	p1, p2, p3 := N/4*1, N/4*2, N/4*3
	chunk0 := make([]string, p1)
	chunk1 := make([]string, p2-p1)
	chunk2 := make([]string, p3-p2)
	chunk3 := make([]string, N-p3)
	copy(chunk0, data[:p1])
	copy(chunk1, data[p1:p2])
	copy(chunk2, data[p2:p3])
	copy(chunk3, data[p3:])

	inps[0] <- chunk0
	inps[1] <- chunk1
	inps[2] <- chunk2
	inps[3] <- chunk3

	sorted := <-res
	for _, line := range sorted {
		writer.WriteString(line)
		writer.WriteString("\n")
	}
}

func mysort(inputPath, outputPath string) {
	fin, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	fout, err := os.Create(outputPath + ".sort")
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

	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})

	for _, line := range data {
		writer.WriteString(line)
		writer.WriteString("\n")
	}
}

func merge(in1, in2, out chan []string) {
	l, r := <-in1, <-in2
	i, j := 0, 0
	m := make([]string, len(l)+len(r))
	for ix := range m {
		switch {
		case i == len(l):
			m[ix] = r[j]
			j++
		case j == len(r):
			m[ix] = l[i]
			i++
		case l[i] > r[j]:
			m[ix] = l[i]
			i++
		default:
			m[ix] = r[j]
			j++
		}
	}

	out <- m
}
