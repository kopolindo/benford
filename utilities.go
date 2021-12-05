package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

var (
	max []int
)

type Pair struct {
	Key   int
	Value float64
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Swap(i, j int) {
	tmp := p[i]
	p[i] = p[j]
	p[j] = tmp
}
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//
//func outputCsv(){
//
//}

func DevStd(s []float64) (devstd float64) {
	xm := Average(s)
	for _, xi := range s {
		devstd += math.Pow((xm - xi), float64(2))
	}
	devstd = math.Round(math.Sqrt((devstd/float64(len(s))))*100) / 100
	return
}

func Average(s []float64) (average float64) {
	for _, i := range s {
		average += i
	}
	average = math.Round((average/float64(len(s)))*100) / 100
	return
}

func Max(s []float64) (max float64) {
	max = s[0]
	for _, val := range s {
		if val >= max {
			max = val
		}
	}
	return
}

func Min(s []float64) (min float64) {
	min = s[0]
	for _, val := range s {
		if val <= min {
			min = val
		}
	}
	return
}

// 1. take the decimal part array and sort it largest -> smallest
// 2. calc difference between total and sum (i.e., N)
// 3. take the first N decimal numbers and add up the respective int part
// note: array of records (couples) {int, dec}

func GenerateParts(perc []float64, tot int) []int {
	var p []Pair
	var intList []int
	sum := 0
	i := 0
	for _, pr := range perc {
		x := pr * float64(tot)
		sum += int(x)
		Int, Frac := math.Modf(x)
		p = append(p, Pair{int(Int), Frac})
		intList = append(intList, int(Int))
		i++
	}
	var pList PairList
	pList = p
	sort.Sort(sort.Reverse(pList))
	delta := tot - sum
	for t := 0; t < delta; t++ {
		max = append(max, pList[t].Key)
	}
	for _, pair := range pList {
		if contains(max, pair.Key) {
			for i, val := range intList {
				if val == pair.Key {
					intList[i]++
				}
			}
		}
	}
	return intList
}

func SSD(input map[int]float64, TOT int) float64 {
	ssd := 0.0
	for i := 0; i < len(input); i++ {
		//fmt.Printf("\"%d\";\"%v\";\"%v\"\n", i+1, input[i+1], thProbs[i])
		ssd += math.Pow((thProbs[i] - input[i+1]), 2)
	}
	return math.Round((10000*ssd)*100) / 100
}

func Compliance(SSD float64) string {
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

func CalcOccurrences(input []int) (out map[int]float64) {
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
