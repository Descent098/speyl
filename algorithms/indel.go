package algorithms

// This file implements the Indel Distance of two strings
//
// # References
//  - https://arxiv.org/abs/2410.09877
//  - https://cs.stackexchange.com/questions/136529/calculating-the-string-similarity-of-an-optimal-alignment
//  - https://cran.r-project.org/web/packages/RapidFuzz/RapidFuzz.pdf
//  - https://arxiv.org/pdf/1909.12877
//  - https://en.wikipedia.org/wiki/Indel

// Calculates the Indel similarity of two strings
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  float32: The indel distance (between 0-1, closer to 1 is more similar)
func IndelSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, IndelDistance)
	return similarity
}

// Calculates the Indel distance of two strings
//
// # Notes
//  - Equivalent to the Levenshtein distance where the cost of substitution is 2, and the cost of insertion or deletion is 1
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The indel distance (edit, delete distance)
func IndelDistance(inputString, targetString string) int {
	// Base cases
	if len(inputString) == 0 {
		return len(targetString)
	}
	if len(targetString) == 0 {
		return len(inputString)
	}

	// If characters match, no cost, move to next characters
	if inputString[0] == targetString[0] {
		return IndelDistance(inputString[1:], targetString[1:])
	}

	// If characters do NOT match, we must perform an operation.
	// We consider two options:
	// 1. Delete inputString[0]: cost 1 + distance of remaining inputString vs targetString
	//    (Effectively, we're removing the current mismatching character from inputString)
	deleteCost := 1 + IndelDistance(inputString[1:], targetString)

	// 2. Insert targetString[0] into inputString: cost 1 + distance of inputString vs remaining targetString
	//    (Effectively, we're adding the current mismatching character from targetString to inputString,
	//     and then we still need to align the rest of inputString)
	insertCost := 1 + IndelDistance(inputString, targetString[1:])

	// Return the minimum of these two options
	return min(deleteCost, insertCost)

}
