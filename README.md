# Benford && CVSS
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
`$ make build`  
It creates a `benford` executable in the same src folder

### Installing
`$ make install`  
It creates a `benford` executable in `$GOPATH/bin/`

### Uninstalling
`$ make uninstall`

### Use
```
$ ./benford -h
  -iterations int
        Number of iterations (default 1)
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
`-verbose` print also additional messages (e.g., compliancy of computed SSD)
`-version` print the version and build of the program

### Examples
#### Example 1
Run 10 times with sample of 1000 scores
```
$ ./benford -sample 1000 -iterations 10
75.88
62.3
82.08
62.26
43.06
64.76
73.14
91.7
69.56
68.82
```

#### Example 2
Run 10 times with sample of 1000 scores, verbose mode
```
$ ./benford -sample 1000 -iterations 10 -verbose
40.66
not good
62.04
not good
74.96
not good
57.08
not good
54.8
not good
84.46
not good
97.64
not good
61.48
not good
77.08
not good
100.54
not even close
```

## Benchmarks

### System Information

- OS: Linux 5.15.4-arch1-1
- Parallel version: 20170922

```$ lscpu | egrep 'Model name|Socket|Thread|NUMA|CPU\(s\)'
CPU(s):                          8
On-line CPU(s) list:             0-7
Model name:                      AMD Ryzen 5 PRO 3500U w/ Radeon Vega Mobile Gfx
Thread(s) per core:              2
Socket(s):                       1
NUMA node(s):                    1
NUMA node0 CPU(s):               0-7
```
Run normal vs parallel
```
10^8
./benford 100000000  12.12s user 1.12s system 106% cpu 12.488 total
parallel ./benford ::: 100000000  10.77s user 1.94s system 109% cpu 11.613 total

10^9
./benford 1000000000  393.67s user 221.88s system 117% cpu 8:43.28 total
parallel ./benford ::: 1000000000  392.62s user 226.67s system 119% cpu 8:37.42 total
```

- OS: Linux 5.15.2-arch1-1
- Parallel version: 20170922

```$ lscpu | egrep 'Model name|Socket|Thread|NUMA|CPU\(s\)'
CPU(s):                          24
On-line CPU(s) list:             0-23
Model name:                      AMD Ryzen 9 3900X 12-Core Processor
Thread(s) per core:              2
Socket(s):                       1
NUMA node(s):                    1
NUMA node0 CPU(s):               0-23
```

Run normal vs parallel
```
10^8
./benford 100000000  6.55s user 0.28s system 105% cpu 6.460 total
parallel ./benford ::: 100000000  6.68s user 0.34s system 105% cpu 6.650 total

10^9
./benford 1000000000  61.77s user 2.32s system 105% cpu 1:00.67 total
parallel ./benford ::: 1000000000  62.11s user 2.28s system 105% cpu 1:00.96 total
```
