# Benford && CVSS
## Index
- [Who](https://github.com/kopolindo/benford#who)
- [Thanks to](https://github.com/kopolindo/benford#thanks-to)
- [Reason](https://github.com/kopolindo/benford#reason)
- [How](https://github.com/kopolindo/benford#how)
- [Expectations](https://github.com/kopolindo/benford#expectations)
- [Program](https://github.com/kopolindo/benford#program)
- [HowTo](https://github.com/kopolindo/benford#howto)
	- [Build](https://github.com/kopolindo/benford#build)
	- [Install](https://github.com/kopolindo/benford#install)
	- [Uninstall](https://github.com/kopolindo/benford#uninstall)
	- [Use](https://github.com/kopolindo/benford#use)
		- [Flag explanation](https://github.com/kopolindo/benford#flag-explanation)
	- [Examples](https://github.com/kopolindo/benford#examples)
		- [Example 1](https://github.com/kopolindo/benford#example-1)
		- [Example 2](https://github.com/kopolindo/benford#example-2)
- [Benchmarks](https://github.com/kopolindo/benford#benchmarks)


## Who
[@kopolindo](https://github.com/kopolindo)  
[@giorgiofox](https://github.com/giorgiofox) (original idea)

## Thanks to
Alex Ely Kossovsky for the awesome statistical review of Chi-Square vs SSD on Benford distributions.

## Reason
We are trying to understand if security assessment results are Benford-like.  
If they are then we can judge security assessment outcomes (big big big numbers must be involved).  

## How
1. Generating CVSS score based on [CVE Details Distribution](https://www.cvedetails.com/cvss-score-distribution.php)
2. Normalizing them with exponential function
3. Calculating SSD (Sum of Squared Deviations). Chi-square is not compatible with Benford distributions, as explained [here](https://www.mdpi.com/2571-905X/4/2/27) by Alex Ely Kossovsky (:beer:)  
4. Iterating this process on and on, to collect more and more data

## Expectations
True, false, blah, not important, hackers gonna hack, just "fun'n'profit" ;)

## Program
Choice: go  
Why: concurrency  
Result: go...home Alex... :(  

## HowTo
### Build
#### Testing race conditions
`$ make race`
It creates a `benford` executable in the same src folder.  
Limitation 8192 concurrent goroutines.  

#### Actual execution
`$ make build`  
It creates a `benford` executable in the same src folder

### Install
`$ make install`  
It creates a `benford` executable in `$GOPATH/bin/`

### Uninstall
`$ make uninstall`

### Use
```
$ ./benford -h
Usage of ./benford:
  -chart
        Create a scattered chart in output folder
  -csv string
        CSV Output filename
  -human
        Human readable vs CSV readable
  -iterations int
        Number of iterations (default 1)
  -max-sample int
        Finish with this sample size (default -1)
  -min-sample int
        Start from this sample size (default -1)
  -sample int
        Size of the sample to be generated
  -verbose
        Verbose, print compliancy
  -version
        Print version
```

#### Flag explanation
`-iterations` is the (int) number of actual runs for the program (default: 1)  
`-sample` is the (int) number of the vulnerabilities among which distribute the scores  
`-min-sample` minimum of the sample set if use case is to range over multiple sample sets. It excludes `-sample`  
`-max-sample` maximum of the sample set if use case is to range over multiple sample sets. It excludes `-sample`  
`-verbose` print also additional messages (e.g., compliancy of computed SSD)  
`-version` print the version and build of the program  
`-chart` generates chart(s) in output folder  
`-human` print in human readable format  
`-csv` output results in output folder, with provided file name

### Examples
#### Example 1
Run 200 times iterations (each iterations returns one SSD)  
Samples spanning from 10000 (vulnerabilities scores) to 20000 (vulnerabilities scores)  
Output:
	- one csv file containing: sample, min, max, average, devstd values
	- one line chart plotting min, max, average (three series overs Y axis) behavior versus sample (X axis)
```
$ ./benford -min-sample 10000 -max-sample 20000 -iterations 200 -chart -csv test.csv
Samples   0% |                                        | (78/10001) [1m9s:2h27m11s]
```
Output:
```
$ ls output
'SSDs result distribution vs samples_line.html'   test.csv

$ head output/test.csv
sample,min,max,average,devstd
10002,52.43,79.79,65.85,4.51
10006,54.42,75.33,65.21,4.11
10003,55.12,77.19,65.32,4.32
10004,57.28,79.51,65.82,4.26
10001,55.52,75.95,65.55,4.18
10005,51.90,78.48,65.46,4.12
10034,52.73,77.86,65.70,4.67
10007,53.30,76.37,65.46,4.63
10008,50.76,75.98,65.41,4.30
```
