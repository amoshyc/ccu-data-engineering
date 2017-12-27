package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"sync"
	"time"
)

var (
	chunkCnt int
	chunkFmt string
)

func main() {
	var inp = flag.String("i", "/dev/stdin", "path/to/input")
	var out = flag.String("o", "/dev/stdout", "path/to/output")
	var fmt = flag.String("f", "/tmp/chunk%05d", "chunkFmt")
	var cnt = flag.Int("c", 1000, "chunkCnt")
	flag.Parse()

	chunkCnt = *cnt
	chunkFmt = *fmt
	externalMergeSort(*inp, *out)
}

func externalMergeSort(inputPath, outputPath string) {
	_ = os.Mkdir(path.Dir(chunkFmt), 0777)

	start := time.Now()
	splitAndSort(inputPath)
	t1 := time.Since(start)

	start = time.Now()
	mergeKWay(outputPath)
	t2 := time.Since(start)

	start = time.Now()
	cleanChunk()
	t3 := time.Since(start)

	fmt.Println("split:", t1)
	fmt.Println("merge:", t2)
	fmt.Println("clean:", t3)
}

func splitAndSort(inputPath string) {
	fin, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	finStat, _ := fin.Stat()
	finSize := finStat.Size()
	chunkBytes := finSize / int64(chunkCnt)
	if finSize%int64(chunkCnt) > 0 {
		chunkBytes++
	}

	scanner := bufio.NewScanner(fin)

	for ix := 0; ix < chunkCnt; ix++ {
		// read data
		chunk := make([]string, 0)
		var nBytes int64
		for nBytes < chunkBytes {
			eof := !scanner.Scan()
			line := scanner.Text()
			chunk = append(chunk, line)
			nBytes += int64(len(line))
			if eof {
				break
			}
		}

		// sort.Slice(chunk, func(i, j int) bool {
		// 	return chunk[i] > chunk[j]
		// })
		chunk = parallelSort(chunk)

		// save chunk
		chunkPath := fmt.Sprintf(chunkFmt, ix)
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			panic(err)
		}
		writer := bufio.NewWriter(chunkFile)
		for _, line := range chunk {
			writer.WriteString(line)
			writer.WriteString("\n")
		}
		writer.Flush()
		chunkFile.Close()

		fmt.Println(chunkPath, " done")
	}
}

func mergeKWay(outputPath string) {
	fout, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	t := newWinnerTree(chunkCnt)
	defer t.Close()

	t.Build(0)

	for !t.Empty() {
		line := t.Pop()
		writer.WriteString(line)
		writer.WriteString("\n")
	}
}

func cleanChunk() {
	for i := 0; i < chunkCnt; i++ {
		chunkPath := fmt.Sprintf(chunkFmt, i)
		err := os.Remove(chunkPath)
		if err != nil {
			panic(err)
		}
	}
}

func parallelSort(data []string) []string {
	res := make([]string, len(data))
	mergeSort(data, res, 0, len(data), 0)
	return res
}

func mergeSort(data []string, res []string, lb, ub int, dep int) {
	if dep >= 5 {
		sort.Slice(data[lb:ub], func(i, j int) bool {
			return data[lb+i] > data[lb+j]
		})
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

type winnerTree struct {
	n    int
	nn   int
	fs   []*os.File
	buf  []*bufio.Scanner
	tree []int
}

func newWinnerTree(n int) *winnerTree {
	w := new(winnerTree)
	w.n = n
	w.nn = 1
	for w.nn < n {
		w.nn <<= 1
	}

	w.tree = make([]int, 2*w.nn-1)
	w.fs = make([]*os.File, n)
	w.buf = make([]*bufio.Scanner, n)
	for ix := range w.buf {
		chunkPath := fmt.Sprintf(chunkFmt, ix)
		var err error
		w.fs[ix], err = os.Open(chunkPath)
		if err != nil {
			panic(err)
		}
		w.buf[ix] = bufio.NewScanner(w.fs[ix])
	}

	return w
}

func (w *winnerTree) Build(u int) {
	if u >= w.nn-1 {
		ix := u - (w.nn - 1)
		if ix >= w.n || !w.buf[ix].Scan() {
			w.tree[u] = -1
		} else {
			w.tree[u] = ix
		}
		return
	}
	lch, rch := 2*u+1, 2*u+2
	w.Build(lch)
	w.Build(rch)
	switch {
	case w.tree[lch] == -1:
		w.tree[u] = w.tree[rch]
	case w.tree[rch] == -1:
		w.tree[u] = w.tree[lch]
	case w.buf[w.tree[lch]].Text() > w.buf[w.tree[rch]].Text():
		w.tree[u] = w.tree[lch]
	default:
		w.tree[u] = w.tree[rch]
	}
}

func (w *winnerTree) Update(u int, val int) {
	if u >= w.nn-1 {
		ix := w.tree[u]
		if ix == -1 {
			return
		}
		if eof := !w.buf[ix].Scan(); eof {
			w.tree[u] = -1
		}
		return
	}
	lch, rch := 2*u+1, 2*u+2
	if w.tree[lch] == val {
		w.Update(lch, val)
	} else {
		w.Update(rch, val)
	}
	switch {
	case w.tree[lch] == -1:
		w.tree[u] = w.tree[rch]
	case w.tree[rch] == -1:
		w.tree[u] = w.tree[lch]
	case w.buf[w.tree[lch]].Text() > w.buf[w.tree[rch]].Text():
		w.tree[u] = w.tree[lch]
	default:
		w.tree[u] = w.tree[rch]
	}
}

func (w *winnerTree) Empty() bool {
	return w.tree[0] == -1
}

func (w *winnerTree) Pop() string {
	ix := w.tree[0]
	line := w.buf[ix].Text()
	w.Update(0, ix)
	return line
}

func (w *winnerTree) Close() {
	for ix := range w.fs {
		w.fs[ix].Close()
	}
}

func (w *winnerTree) Pr(u int, ind int) {
	if u >= 2*w.nn-1 {
		return
	}

	ix := w.tree[u]
	for i := 0; i < ind; i++ {
		fmt.Print("   ")
	}
	if ix == -1 {
		fmt.Printf("%3d (xxx)\n", -1)
	} else {
		fmt.Printf("%3d (%3s)\n", ix, w.buf[ix].Text())
	}

	w.Pr(2*u+1, ind+1)
	w.Pr(2*u+2, ind+1)
}

