package algorithms

// Calculates the Levenshtein similarity of two strings
//
// # Parameters
//
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  float32: The similarity (between 0-1, closer to 1 is more similar)
func LevenshteinSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, LevenshteinDistance)
	return similarity
}

// Calculates the Levenshtein distance of two strings
//
// # Notes
//  - This solution utilizes the dynamic programming approach, not the recursive one
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenshtein distance (add, edit, delete distance)
func LevenshteinDistance(inputString, targetString string) int {
	return DynamicLevenshtein(inputString, targetString)
}

// Calculates the Levenshtein distance of two strings recursively
//
// # Notes
//  - Heavily inspired by the recursive haskel implementation on wikipedia https://en.wikipedia.org/wiki/Levenshtein_distance#Recursive
//  - Very slow, roughly O(3^n)
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenshtein distance (add, edit, delete distance)
func RecursiveLevenshtein(inputString, targetString string) int {
	if len(inputString) == 0 {
		return len(targetString)
	}
	if len(targetString) == 0 {
		return len(inputString)
	}

	firstInputChar := inputString[0]
	firstTargetChar := targetString[0]
	restInputString := inputString[1:]
	restTargetString := targetString[1:]

	if firstInputChar == firstTargetChar {
		return RecursiveLevenshtein(restInputString, restTargetString)
	}

	return 1 + min(
		RecursiveLevenshtein(inputString, restTargetString),     // Add
		RecursiveLevenshtein(restInputString, targetString),     // Delete
		RecursiveLevenshtein(restInputString, restTargetString), // Edit/replace
	)

}

func RecursiveLevenshteinSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, RecursiveLevenshtein)
	return similarity
}

// A dynamic-programming based implementation of Levenshtein distance
//
// # Notes
//  - Relies on Wagner–Fischer algorithm https://en.wikipedia.org/wiki/Wagner%E2%80%93Fischer_algorithm#Calculating_distance
//  - More details: https://gist.github.com/Descent098/401c2ca6bdf3fa655738e7a1ddf1aeee
//  - Faster than the recursive solution, runs in roughly O(m*n) where m and n are the size of strings
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Levenshtein distance (add, edit, delete distance)
func DynamicLevenshtein(inputString, targetString string) int {
	// Convert to runes to avoid weird encoding issues
	inputStringRunes := []rune(inputString)
	targetStringRunes := []rune(targetString)

	inputStringLength := len(inputStringRunes)
	targetStringLength := len(targetStringRunes)

	// Create a 2D matrix
	matrix := make([][]int, inputStringLength+1)
	for i := range matrix {
		matrix[i] = make([]int, targetStringLength+1)
	}

	// Initialize base cases
	for i := 0; i <= inputStringLength; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= targetStringLength; j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= inputStringLength; i++ {
		for j := 1; j <= targetStringLength; j++ {
			if inputStringRunes[i-1] == targetStringRunes[j-1] {
				// Characters match, no cost added
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = 1 + min(
					matrix[i][j-1],   // Add
					matrix[i-1][j],   // Delete
					matrix[i-1][j-1], // Edit/replace
				)
			}
		}
	}

	return matrix[inputStringLength][targetStringLength]
}

// A recursive Levenshtein distance using the Damerau–Levenshtein distance
//
// # Notes
//  - Relies on Damerau–Levenshtein distance, which is the Levenshtein distance + transpositions
//  - More details: https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance
//  - Relies on memoization for performance and accuracy: https://en.wikipedia.org/wiki/Memoization
//
// # Parameters
//  inputString (string): The first string to use for the comparison
//  targetString (string): The second string to use for the comparison
//
// # Returns
//  int: The Damerau–Levenshtein distance (add, edit, delete, transpose distance)
func DamerauLevenshtein(input, target string) int {
	// Create a memoization cache
	cache := make(map[string]int)

	// Define an inner function to memoize
	var calculateDamerauLevenshteinDistance func(a, b string) int

	calculateDamerauLevenshteinDistance = func(a, b string) int {
		// Handle simple cases
		key := a + "|" + b
		if val, exists := cache[key]; exists {
			return val
		}

		// Empty strings
		if len(a) == 0 {
			cache[key] = len(b)
			return cache[key]
		}
		if len(b) == 0 {
			cache[key] = len(a)
			return cache[key]
		}

		// Exact matches for first char, go to next char
		if a[0] == b[0] {
			cache[key] = calculateDamerauLevenshteinDistance(a[1:], b[1:])
			return cache[key]
		}

		// Calculate Levenshtein Distance
		insert := calculateDamerauLevenshteinDistance(a, b[1:])
		delete := calculateDamerauLevenshteinDistance(a[1:], b)
		replace := calculateDamerauLevenshteinDistance(a[1:], b[1:])
		minCost := 1 + min(insert, delete, replace)

		// Check for transposition to add Damerau changes
		if len(a) > 1 &&
			len(b) > 1 &&
			a[0] == b[1] &&
			a[1] == b[0] {
			transpose := 1 + calculateDamerauLevenshteinDistance(a[2:], b[2:])
			minCost = min(minCost, transpose)
		}

		// Update memoize cache
		cache[key] = minCost
		return minCost
	}

	return calculateDamerauLevenshteinDistance(input, target)
}

func DamerauLevenshteinSimilarity(inputString, targetString string) float32 {
	similarity := CalculateSimilarity(inputString, targetString, DamerauLevenshtein)
	return similarity
}
