# Speyl

A Golang Based spell checker

## Installation

Add:

```go
import "github.com/Descent098/speyl"
```
 
To your go file, then run `go mod tidy` to download. Here's a minimal example:

```go
package main

import (
	"fmt"

	"github.com/Descent098/speyl"
)

func main() {
	validWords := []string{"hi", "hello", "bonjour", "alumni"}

	inputWord := "alumni"
	s := speyl.SuggestWord(inputWord, validWords) // Returns algorithms.Suggestion{Likelihood: 1.0, Word: "alumni"}
	// Prints: alumni is alumni with %100.00 Likelihood
	fmt.Printf("%s is %s with %%%.2f Likelihood\n", inputWord, s.Word, s.Likelihood*100)

	misspeltWord := "almni"
	s = speyl.SuggestWord(misspeltWord, validWords) // Returns algorithms.Suggestion{Likelihood: 0.944, Word: "alumni"}
	// Prints: almni is alumni with %94.44 Likelihood
	fmt.Printf("%s is %s with %%%.2f Likelihood\n", misspeltWord, s.Word, s.Likelihood*100)
}
```


## Usage

The generally recommended usage is to use the `SuggestWord()` function in the main package, this uses the [Jaro Similarity](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance), and is the fastest algorithm by far:

```go 
func SuggestWord(word string, validWords []string) algorithms.Suggestion{}
```

Which returns a struct:

```go
type Suggestion struct {
	Likelihood float32 // How confident the suggestion is
	Word       string  // The suggested word
}
```

If you want to use a different algorithm you can do so with:

```go
func SuggestWordWithSpecificAlgorithm(word string, validWords []string, algorithm algorithms.SimilarityAlgorithm) algorithms.Suggestion {}
```

You can find the various available `SimilarityAlgorithm`'s in the `algorithm` package.

## Performance

Below is the performance tests of the various algorithms and their implementations. They were tested using `words.txt` a corpus of ~370,000 words. There were two separate tests. The first was the synchronus execution using `algorithms.SuggestWord()`. 

| Algorithm | Time Taken (miliseconds) |
|-----------|--------------------------|
| Jaro | 27 |
| Levenshtein (Dynamic Programming) | 105 |
| Indel | 1,975 |
| Levenshtein (Recursive Damerau) | 2,997 | 
| Levenshtein (Recursive) | 18,078 |

You may assume that `SuggestWord()` is concurrent, but actually because of the execution speed of the Jaro and Dynamic Programming Levenshtein, it was slower to do the task asynchronously than synchronously. Since these would be the two that are most likely to get practical use I left the main implementations synchronous. If your list of valid words is over ~1 mil it might be worth implementing the concurrent version yourself. Here the code I tried (I also tried with channels and it was still slower):

```go
func SuggestWord(word string, validWords []string, algorithm algorithms.SimilarityAlgorithm) algorithms.Suggestion {
	writeResultLock := sync.Mutex{}

	highestRatio := float32(0)
	currentSuggestion := ""

	wg := sync.WaitGroup{}

	for _, currentWord := range validWords {
		wg.Add(1)
		go func(wordToCompare string) {
			defer wg.Done()
			defer writeResultLock.Unlock()
			res := algorithm(word, wordToCompare)
			writeResultLock.Lock()
			if res > highestRatio {
				highestRatio = res
				currentSuggestion = currentWord
			}
		}(currentWord)

	}
	wg.Wait()
	return algorithms.Suggestion{
		Likelihood: highestRatio,
		Word:       currentSuggestion,
	}
}
```

| Algorithm | Time Taken (miliseconds) |
|-----------|--------------------------|
| Jaro | 70 |
| Levenshtein (Dynamic Programming) | 163 |
| Indel | 205 |
| Levenshtein (Recursive Damerau) | 620 | 
| Levenshtein (Recursive) | 1542 |
