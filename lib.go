// A library to help calculate the similarity between two strings
//
// # Functions
//
//	CheckSimilarity(): Check the similarity of a word to a corpus of validWords
//	LoadWords(): Returns a slice of a huge number of words (>350,000)
package speyl

import (
	_ "embed"
	"strings"
	"sync"

	"github.com/Descent098/speyl/algorithms"
)

//go:embed words.txt
var wordsContent string

// Helper function to load a default corpus of over 350,000 words
//
// # Returns
//
//	[]string: A slice with the words in the corpus
func LoadWords() []string {
	result := strings.Split(wordsContent, "\r\n")
	return result
}

// Utility function that does the heavy lifting, essentially a general solution that will take in any SimilarityAlgorithm
//
// # Parameters
//
//	inputWord (*C.char): The word to find a similar word for
//	algorithm (algorithms.SimilarityAlgorithm): The algorithm to use to calculate the similarity of the words
//	validWords ([]string): A slice with the words in the corpus
//
// # Returns
//
//	Suggestion: A pointer to the suggestion struct with the word and it's likelihood
func CheckSimilarity(word string, algorithm algorithms.SimilarityAlgorithm, validWords []string) algorithms.Suggestion {
	result := algorithms.SuggestWord(word, validWords, algorithm)
	return result
}

func SuggestWord(word string, validWords []string) algorithms.Suggestion {
	writeResultLock := sync.Mutex{}

	highestRatio := float32(0)
	currentSuggestion := ""

	wg := sync.WaitGroup{}

	for _, currentWord := range validWords {
		wg.Add(1)
		go func(wordToCompare string) {
			defer wg.Done()
			defer writeResultLock.Unlock()
			res := algorithms.JaroSimilarity(word, wordToCompare)
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
