# Benford && CVSS
## Reason
We are trying to understand if security assessment results are Benford-like.  
If they are then we can judge security assessment outcomes (big big big numbers must be involved).  

## How
1. Generating CVSS score based on [CVE Details Distribution](https://www.cvedetails.com/cvss-score-distribution.php)
2. Normalizing them with exponential function
3. Calculating SSD (Sum of Squared Deviations). Chi-square is not compatible with Benford distributions, as explained [here](https://www.mdpi.com/2571-905X/4/2/27) by Alex Ely Kossovsky (:beer:)  
4. Iterating this process on and on, to collect more and more data

## Expectations
True, false, blah, not important, hack wanna hack, just "fun'n'profit" ;)

## Program
Choice: go
Why: concurrency
Result: go...home Alex... :(
