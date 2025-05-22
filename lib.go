// A library to help calculate the similarity between two strings
//
// # Functions
//
//	CheckSimilarity(): Check the similarity of a word to a corpus of validWords
//	LoadWords(): Returns a slice of a huge number of words (>350,000)
package speyl

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Descent098/speyl/algorithms"
)

// Gets the file path of the currently running go file
//
// # Returns
//
//	string: the path to the current file in go
func getCurrentFilePath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get caller information")
	}
	return filename
}

// Helper function to load a default corpus of over 350,000 words
//
// # Returns
//
//	[]string: A slice with the words in the corpus
func LoadPremadeWords() []string {
	currentFileDir := filepath.Dir(getCurrentFilePath())
	filePath := filepath.Join(currentFileDir, "words.txt")

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	result := strings.Split(string(content), "\r\n")
	return result
}

// Used to get a suggested word with a specific algorithm
//
// # Parameters
//
//	inputWord (string): The word to find a similar word for
//	validWords ([]string): A slice with the words in the corpus
//	algorithm (algorithms.SimilarityAlgorithm): The algorithm to use to calculate the similarity of the words
//
// # Returns
//
//	algorithms.Suggestion: A pointer to the suggestion struct with the word and it's likelihood
func SuggestWordWithSpecificAlgorithm(word string, validWords []string, algorithm algorithms.SimilarityAlgorithm) algorithms.Suggestion {
	result := algorithms.SuggestWord(word, validWords, algorithm)
	return result
}

// Used to get a suggestion using Jaro Similarity, fastest solution available
//
// # Parameters
//
//	inputWord (string): The word to find a similar word for
//	validWords ([]string): A slice with the words that are considered valid
//
// # Returns
//
//	Suggestion: A suggestion struct with the word and it's likelihood
func SuggestWord(word string, validWords []string) algorithms.Suggestion {
	highestRatio := float32(0)
	currentSuggestion := ""

	for _, currentWord := range validWords {
		res := algorithms.JaroSimilarity(word, currentWord)
		if res > highestRatio {
			highestRatio = res
			currentSuggestion = currentWord
		}
	}
	return algorithms.Suggestion{
		Likelihood: highestRatio,
		Word:       currentSuggestion,
	}
}
