package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func randString() string {
	lb, ub := 0x4E00, 0x9FA5
	res := make([]rune, 4)
	for ix := range res {
		res[ix] = rune(rand.Intn(ub-lb+1) + lb)
	}
	return string(res)
}

func generate(outputPath string, n int) {
	fout, err := os.Create(outputPath + "." + strconv.Itoa(n))
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	for i := 0; i < (1 << uint(n)); i++ {
		writer.WriteString(string(rand.Intn(1000000)))
		writer.WriteString("\n")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	out := flag.String("o", "/tmp/amoshyc/in", "output location")
	flag.Parse()

	for i := 20; i < 30; i += 2 {
		fmt.Print(i, ": ")
		startTime := time.Now()
		generate(*out, i)
		fmt.Println(time.Since(startTime))
	}
}
