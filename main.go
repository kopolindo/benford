package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	charts "benford/charts"
	"benford/structure"
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
	mainMutex        sync.Mutex
)

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
func worker(sample int, iterations int, mainWg *sync.WaitGroup, resChan chan structure.Result) {
	defer mainWg.Done()
	var res structure.Result
	var ssdResults []float64
	var wg sync.WaitGroup
	mainMutex.Lock()
	wg.Add(iterations)
	SSDsResultsChan := make(chan float64, 1)
	go func() {
		// Close channel
		wg.Wait()
		close(SSDsResultsChan)
	}()
	for it := 0; it < iterations; it++ {
		// Create channel to make goroutine and main routine communicate
		go func(sample int, localWg *sync.WaitGroup, result chan float64) {
			defer localWg.Done()
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
		}(sample, &wg, SSDsResultsChan)
	}
	// Fetch value from gorouting
	for r := range SSDsResultsChan {
		ssdResults = append(ssdResults, r)
	}
	res.Lock()
	res.SSDs = ssdResults
	res.Sample = sample
	res.Average = Average(ssdResults)
	res.Max = Max(ssdResults)
	res.Min = Min(ssdResults)
	res.DevStd = DevStd(ssdResults)
	res.Unlock()
	resChan <- res
	mainMutex.Unlock()
}

func main() {
	resultChannel := make(chan structure.Result, 1)
	var m sync.Mutex
	var mainWg sync.WaitGroup
	sampleSetSize := maxSample - minSample + 1
	mainWg.Add(sampleSetSize)
	if verbose {
		fmt.Printf("minSample: %d\nmaxSample: %d\nsampleSetSize: %d\n",
			minSample, maxSample, sampleSetSize)
	}
	go func() {
		mainWg.Wait()
		close(resultChannel)
	}()
	for sample = minSample; sample <= maxSample; sample++ {
		go worker(sample, iterations, &mainWg, resultChannel)
	}
	for workerResult := range resultChannel {
		if chart {
			m.Lock()
			var scatterChart charts.ScatterData
			scatterChart.Create(workerResult)
			m.Unlock()
		}
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
}
