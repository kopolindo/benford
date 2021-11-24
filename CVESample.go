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

var realCVSSDistribution = [4]riskDistribution{
	{"C", 9.0, 10.0, 0.115},
	{"H", 7.0, 8.9, 0.207},
	{"M", 4.0, 6.9, 0.574},
	{"L", 0.0, 3.9, 0.104},
}

func firstLeftDigit(a float64) int {
	if a < 10.0 {
		return int(a)
	}
	a = float64(firstLeftDigit(a / 10.0))
	return int(a)
}

func GenerateFirstDigitCVSSScores(tot int) []int {
	var FirstDigitCVSSScores []int
	for _, risk := range realCVSSDistribution {
		//fmt.Println("Printing", risk.level, "scores (", risk.minScore, ",", risk.maxScore, ")")
		for i := 0; i < int(risk.percentage*float64(tot)); i++ {
			fld := firstLeftDigit(math.Exp(risk.minScore + rand.Float64()*(risk.maxScore-risk.minScore)))
			FirstDigitCVSSScores = append(FirstDigitCVSSScores, fld)
		}
	}
	return FirstDigitCVSSScores
}
