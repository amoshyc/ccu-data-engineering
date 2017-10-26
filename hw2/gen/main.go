package main

import (
	"flag"
	"fmt"
	"math/rand"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	n := flag.Int("n", 1000, "number of data")
	flag.Parse()

	for ix := 0; ix < *n; ix++ {
		fmt.Println(rand.Intn(1000000))
		// fmt.Println(randString())
	}
}
