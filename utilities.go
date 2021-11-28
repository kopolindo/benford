package main

import (
	"math"
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
