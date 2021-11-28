package main

import (
	"math"
	"math/rand"
	"sync"
)

type riskDistribution struct {
	level      string
	minScore   float64
	maxScore   float64
	percentage float64
}

type IntSlice []int

var (
	CVEPercentageDistribution []float64
	realCVSSDistribution      = [4]riskDistribution{
		{"C", 9.0, 10.0, 0.115},
		{"H", 7.0, 8.9, 0.207},
		{"M", 4.0, 6.9, 0.574},
		{"L", 0.0, 3.9, 0.104},
	}
	wg sync.WaitGroup
)

func firstLeftDigit(a float64) int {
	if a < 10.0 {
		return int(a)
	}
	a = float64(firstLeftDigit(a / 10.0))
	return int(a)
}

func RiskCalc(risk riskDistribution, quantities int, workg *sync.WaitGroup, result chan []int) {
	var out []int
	defer workg.Done()
	for k := 0; k < quantities; k++ {
		score := math.Round((risk.minScore+rand.Float64()*(risk.maxScore-risk.minScore))*10) / 10
		fld := firstLeftDigit(math.Exp(score))
		out = append(out, fld)
	}
	result <- out
	return
}

func GenerateFirstDigitCVSSScores(tot int) []int {
	var FirstDigitCVSSScores []int
	for _, r := range realCVSSDistribution {
		CVEPercentageDistribution = append(CVEPercentageDistribution, r.percentage)
	}
	quantities := GenerateParts(CVEPercentageDistribution, tot)
	var wg sync.WaitGroup
	wg.Add(len(realCVSSDistribution))
	for i, risk := range realCVSSDistribution {
		tmp := make(chan []int, 1)
		go RiskCalc(risk, quantities[i], &wg, tmp)
		res := <-tmp
		close(tmp)
		//for k := 0; k < quantities[i]; k++ {
		//	score := math.Round((risk.minScore+rand.Float64()*(risk.maxScore-risk.minScore))*10) / 10
		//	//fmt.Println(score)
		//	fld := firstLeftDigit(math.Exp(score))
		//}
		FirstDigitCVSSScores = append(FirstDigitCVSSScores, res...)
	}
	wg.Wait()
	return FirstDigitCVSSScores
}
