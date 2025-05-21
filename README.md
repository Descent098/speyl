# Speyl

A Golang Based spell checker


## Performance

Below is the performance tests of the various algorithms and their implementations. They were tested using `words.txt` a corpus of ~370,000 words. There were two separate tests. The first was the synchronus execution using `algorithms.SuggestWord()`. 

| Algorithm | Time Taken (miliseconds) |
|-----------|--------------------------|
| Jaro | 27 |
| Levenshtein (Dynamic Programming) | 105 |
| Levenshtein (Recursive) | 18078 |
| Indel | 1975 |
| Levenshtein (Recursive Damerau) | 2997 | 

The second was using the concurrent model:

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
