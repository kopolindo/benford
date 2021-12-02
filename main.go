package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
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
	verbose          bool
	mainWG           sync.WaitGroup
	ssdResults       []float64
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations")
	flag.IntVar(&sample, "sample", 0, "Size of the sample to be generated")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.BoolVar(&verbose, "verbose", false, "Verbose, print compliancy")
	flag.Parse()
}

func main() {
	// Print version
	if version {
		fmt.Println("Version:\t", Version)
		fmt.Println("Build:\t\t", BuildCommitShort)
		os.Exit(0)
	}
	// If no flag is provided
	if sample < 1 || iterations < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// If sample is not is not provided
	if sample < 1 {
		fmt.Println("Sample must be at least 1.\nThe greater the better.\nFrom great samples come great statistics.\nUse -sample flag to provide sample.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Create working group (one for every iteration)
	mainWG.Add(iterations)
	// Initialize SSD result array
	for it := 0; it < iterations; it++ {
		// Create channel to make goroutine and main routine communicate
		SSDs := make(chan float64, 1)
		go func(sample int, workg *sync.WaitGroup, result chan float64) {
			defer workg.Done()
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
		}(sample, &mainWG, SSDs)
		// Fetch value from gorouting
		//ssd := <-SSDs
		ssdResults = append(ssdResults, <-SSDs)
		// Close channel
		close(SSDs)
	}
	fmt.Println("Min:", Min(ssdResults))
	fmt.Println("Max:", Max(ssdResults))
	fmt.Println("Average:", Average(ssdResults))

	// Print the output
	//fmt.Println(ssd)
	// If verbose print also the compliancy
	//if verbose {
	//	fmt.Println(Compliance(ssd))
	//}
}
