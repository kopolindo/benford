package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	charts "benford/charts"
)

const (
	DEFAULT_MIN = 100
	DEFAULT_MAX = 100000
)

var (
	thProbs          = [9]float64{0.301, 0.176, 0.125, 0.097, 0.079, 0.067, 0.058, 0.057, 0.046}
	test             = [9]float64{0.327, 0.149, 0.107, 0.079, 0.076, 0.082, 0.061, 0.055, 0.064}
	sample           int
	minSample        int
	maxSample        int
	iterations       int
	Version          string
	BuildCommitShort string
	version          bool
	verbose          bool
	humanReadable    bool
	chart            bool
	iterationsWG     sync.WaitGroup
	sampleWG         sync.WaitGroup
	m                sync.RWMutex
)

type Result struct {
	Sample  int
	Average float64
	Min     float64
	Max     float64
	DevStd  float64
	SSDs    []float64
}

func checkConditions() {
	// Print version
	if version {
		fmt.Println("Version:\t", Version)
		fmt.Println("Build:\t\t", BuildCommitShort)
		os.Exit(0)
	}
	// If no sample flag is provided
	if !IsFlagPassed("sample") &&
		!IsFlagPassed("min-sample") &&
		!IsFlagPassed("max-sample") {
		fmt.Println("You need to specify at least one sample ;)\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if IsFlagPassed("sample") &&
		(IsFlagPassed("min-sample") || IsFlagPassed("max-sample")) {
		fmt.Println("The following flags are incompatible one with each other:")
		fmt.Println("\tsample with min-sample and max-sample\n")
		fmt.Println("You use sample for a one-shot execution")
		fmt.Println("You use min-sample and max-sample to range samples\n")
		fmt.Println("eg: benford -iterations 1000 -min-sample 100 -max-sample 1000\n")
		fmt.Println("This would execute the program 900 times, using increasing samples")
		os.Exit(1)
	}
	// If minSample and maxSample are set, then sample is not needed
	if IsFlagPassed("sample") {
		minSample = sample
		maxSample = sample
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations")
	flag.IntVar(&sample, "sample", 0, "Size of the sample to be generated")
	flag.IntVar(&minSample, "min-sample", -1, "Start from this sample size")
	flag.IntVar(&maxSample, "max-sample", -1, "Finish with this sample size")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.BoolVar(&verbose, "verbose", false, "Verbose, print compliancy")
	flag.BoolVar(&humanReadable, "human", false, "Human readable vs CSV readable")
	flag.BoolVar(&chart, "chart", false, "Create a scattered chart in output folder")
	flag.Parse()
	checkConditions()
}

// Worker function (for sample)
func worker(sample int, iterations int, mainWg *sync.WaitGroup, resChan chan Result) {
	defer mainWg.Done()
	var res Result
	// Initialize SSD result array
	var ssdResults []float64
	SSDsResultsChan := make(chan float64, 1)
	go func(result chan float64) {
		// Fetch value from gorouting
		for r := range result {
			ssdResults = append(ssdResults, r)
		}
	}(SSDsResultsChan)
	var wg sync.WaitGroup
	wg.Add(iterations)
	for it := 0; it < iterations; it++ {
		// Create channel to make goroutine and main routine communicate
		go func(sample int, result chan float64) {
			defer wg.Done()
			var keys []int
			// Generate CVSS scores, normalize them (Exp) and take the first digit
			fdCVSSScores := GenerateFirstDigitCVSSScores(sample)
			// count occurrences of first left digits
			occurrences := CalcOccurrences(fdCVSSScores)
			// Here the order part
			for k := range occurrences {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			// Communicate with main routine
			result <- SSD(occurrences, sample)
		}(sample, SSDsResultsChan)
	}
	// Close channel
	wg.Wait()
	close(SSDsResultsChan)
	m.Lock()
	res.SSDs = ssdResults
	res.Sample = sample
	res.Average = Average(ssdResults)
	res.Max = Max(ssdResults)
	res.Min = Min(ssdResults)
	res.DevStd = DevStd(ssdResults)
	m.Unlock()
	resChan <- res
}

func main() {
	sampleSetSize := maxSample - minSample + 1
	fmt.Printf("minSample: %d\nmaxSample: %d\nsampleSetSize: %d\n",
		minSample, maxSample, sampleSetSize)
	resultChannel := make(chan Result, 1)
	var workerResult Result
	go func(ch chan Result) {
		for i := range ch {
			workerResult = i
			if humanReadable {
				fmt.Println("Min:", workerResult.Min)
				fmt.Println("Max:", workerResult.Max)
				fmt.Println("Average:", workerResult.Average)
				fmt.Println("DevStd", workerResult.DevStd)
			} else {
				fmt.Printf("%d;%.2f;%.2f;%.2f;%.2f\n",
					workerResult.Sample,
					workerResult.Min,
					workerResult.Max,
					workerResult.Average,
					workerResult.DevStd)
			}
		}
	}(resultChannel)
	var mainWg sync.WaitGroup
	mainWg.Add(sampleSetSize)
	for sample = minSample; sample <= maxSample; sample++ {
		go worker(sample, iterations, &mainWg, resultChannel)
		if chart {
			var scatterChart charts.ScatterData
			scatterChart.Create(strconv.Itoa(sample), workerResult.SSDs)
		}
	}
	mainWg.Wait()
	close(resultChannel)
}
