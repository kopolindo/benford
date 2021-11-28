package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

var thProbs = [9]float64{0.301, 0.176, 0.125, 0.097, 0.079, 0.067, 0.058, 0.057, 0.046}
var test = [9]float64{0.327, 0.149, 0.107, 0.079, 0.076, 0.082, 0.061, 0.055, 0.064}

func compliance(SSD float64) string {
	result := "default"
	if SSD < 0 {
		fmt.Println("ssd cannot be negative :/")
		os.Exit(1)
	}
	switch ssd := SSD; {
	case ssd < 2.0:
		result = "perfect"
	case ssd < 25.0:
		result = "meh"
	case ssd < 100.0:
		result = "not good"
	case ssd >= 100.0:
		result = "not even close"
	default:
		result = "default"
	}
	return result
}

func SSD(input [9]float64, TOT int) float64 {
	ssd := 0.0
	for i, p := range input {
		fmt.Printf("\"%d\";\"%v\";\"%v\"\n", i+1, p, thProbs[i])
		ssd += math.Pow(float64(TOT)*(p-thProbs[i]), 2)
	}
	return ssd
}

func calcOccurrences(input []int) (out map[int]float64) {
	tot := 0
	tmp := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
		7: 0,
		8: 0,
		9: 0,
	}
	out = make(map[int]float64)
	for _, i := range input {
		tmp[i]++
	}
	for _, v := range tmp {
		tot += v
	}
	for k, v := range tmp {
		out[k] = float64(v) / float64(tot)
	}
	return
}

func main() {
	var keys []int
	TOT, _ := strconv.Atoi(os.Args[1])
	rand.Seed(time.Now().UnixNano())
	fdCVSSScores := GenerateFirstDigitCVSSScores(TOT)
	occurrences := calcOccurrences(fdCVSSScores)
	for k := range occurrences {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println(k, occurrences[k])
	}
	Hist(occurrences)
	//ssd := SSD(occurrences, TOT)
	//fmt.Println(compliance(ssd))
}
