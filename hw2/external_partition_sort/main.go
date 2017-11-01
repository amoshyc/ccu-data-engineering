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
	externalPartitionSort(*inp, *out)
}

func externalPartitionSort(inputPath, outputPath string) {
	_ = os.Mkdir(path.Dir(chunkFmt), 0777)

	start := time.Now()
	partition(inputPath)
	t1 := time.Since(start)

	start = time.Now()
	sortAndConcat(outputPath)
	t2 := time.Since(start)

	start = time.Now()
	cleanChunk()
	t3 := time.Since(start)

	fmt.Println("parti:", t1)
	fmt.Println("cncat:", t2)
	fmt.Println("clean:", t3)
}

func partition(inputPath string) {
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

	pivots := make([]string, chunkCnt-1)
	pivotScanner := bufio.NewScanner(fin)
	for ix := range pivots {
		var nBytes int64
		var line string
		for nBytes < chunkBytes {
			eof := !pivotScanner.Scan()
			line = pivotScanner.Text()
			nBytes += int64(len(line))
			if eof {
				break
			}
		}
		pivots[ix] = line
	}

	sort.Slice(pivots, func(i, j int) bool {
		return pivots[i] > pivots[j]
	})

	fmt.Println(len(pivots))

	files := make([]*bufio.Writer, chunkCnt)
	for ix := range files {
		chunkPath := fmt.Sprintf(chunkFmt, ix)
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			panic(err)
		}
		defer chunkFile.Close()
		files[ix] = bufio.NewWriter(chunkFile)
		defer files[ix].Flush()
	}

	fin.Seek(0, 0)
	dataScanner := bufio.NewScanner(fin)
	for dataScanner.Scan() {
		line := dataScanner.Text()
		ix := sort.Search(len(pivots), func(i int) bool { // binary search
			return pivots[i] <= line
		})

		files[ix].WriteString(line)
		files[ix].WriteString("\n")
	}
}

func sortAndConcat(outputPath string) {
	fout, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	for ix := 0; ix < chunkCnt; ix++ {
		chunkPath := fmt.Sprintf(chunkFmt, ix)
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(chunkFile)

		data := make([]string, 0)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}

		if len(data) > 0 {
			// sort.Slice(data, func(i, j int) bool {
			// 	return data[i] > data[j]
			// })
			data = parallelSort(data)

			for _, line := range data {
				writer.WriteString(line)
				writer.WriteString("\n")
			}
		}

		chunkFile.Close()
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
