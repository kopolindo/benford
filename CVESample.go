package main

import (
	"math"
	"math/rand"
)

type riskDistribution struct {
	level      string
	minScore   float64
	maxScore   float64
	percentage float64
}

var (
	CVEPercentageDistribution []float64
	realCVSSDistribution      = [4]riskDistribution{
		{"C", 9.0, 10.0, 0.115},
		{"H", 7.0, 8.9, 0.207},
		{"M", 4.0, 6.9, 0.574},
		{"L", 0.0, 3.9, 0.104},
	}
)

func firstLeftDigit(a float64) int {
	if a < 10.0 {
		return int(a)
	}
	a = float64(firstLeftDigit(a / 10.0))
	return int(a)
}

func GenerateFirstDigitCVSSScores(tot int) []int {
	var FirstDigitCVSSScores []int
	for _, r := range realCVSSDistribution {
		CVEPercentageDistribution = append(CVEPercentageDistribution, r.percentage)
	}
	quantities := GenerateParts(CVEPercentageDistribution, tot)
	for i, risk := range realCVSSDistribution {
		for k := 0; k < quantities[i]; k++ {
			score := math.Round((risk.minScore+rand.Float64()*(risk.maxScore-risk.minScore))*10) / 10
			//fmt.Println(score)
			fld := firstLeftDigit(math.Exp(score))
			FirstDigitCVSSScores = append(FirstDigitCVSSScores, fld)
		}
	}
	return FirstDigitCVSSScores
}
