package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

var (
	thProbs          = [9]float64{0.301, 0.176, 0.125, 0.097, 0.079, 0.067, 0.058, 0.057, 0.046}
	test             = [9]float64{0.327, 0.149, 0.107, 0.079, 0.076, 0.082, 0.061, 0.055, 0.064}
	sample           int
	iterations       int
	Version          string
	BuildCommitShort string
	version          bool
	keys             []int
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations")
	flag.IntVar(&sample, "sample", 1, "Size of the sample to be generated")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Parse()
}

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
		result = "quite good"
	case ssd < 100.0:
		result = "not good"
	case ssd >= 100.0:
		result = "not even close"
	default:
		result = "default"
	}
	return result
}

func SSD(input map[int]float64, TOT int) float64 {
	ssd := 0.0
	for i := 0; i < len(input); i++ {
		//fmt.Printf("\"%d\";\"%v\";\"%v\"\n", i+1, input[i+1], thProbs[i])
		ssd += math.Pow((thProbs[i] - input[i+1]), 2)
	}
	return math.Round((10000*ssd)*100) / 100
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
	if version {
		fmt.Println("Version:\t", Version)
		fmt.Println("Build:\t\t", BuildCommitShort)
		os.Exit(0)
	}
	for it := 0; it < iterations; it++ {
		fdCVSSScores := GenerateFirstDigitCVSSScores(sample)
		occurrences := calcOccurrences(fdCVSSScores)
		for k := range occurrences {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		/*
			for _, k := range keys {
				fmt.Println(k, occurrences[k])
			}*/
		ssd := SSD(occurrences, sample)
		fmt.Println(ssd)
		//fmt.Println(compliance(ssd))
	}
}
