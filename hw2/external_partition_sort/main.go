package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
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
	in := make(chan []string, 1)
	out := make(chan []string, 1)
	go mergesort(in, out, 0)
	in <- data
	return <-out
}

func mergesort(in chan []string, out chan []string, dep int) {
	data := <-in

	if dep >= 5 {
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		out <- data
		return
	}

	N := len(data)
	in1 := make(chan []string, 1)
	in2 := make(chan []string, 1)
	res1 := make(chan []string, 1)
	res2 := make(chan []string, 1)
	go mergesort(in1, res1, dep+1)
	go mergesort(in2, res2, dep+1)
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
