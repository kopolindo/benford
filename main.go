package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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

func SSD(input [9]float64) float64 {
	ssd := 0.0
	for i, p := range input {
		ssd += math.Pow((p - thProbs[i]), 2)
	}
	return ssd
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//fmt.Println("Benford")
	/*ssd := SSD(test)
	fmt.Println(ssd)
	fmt.Println(compliance(ssd))*/
	fdCVSSScores := GenerateFirstDigitCVSSScores(1500)
	fmt.Println(len(fdCVSSScores))
}
