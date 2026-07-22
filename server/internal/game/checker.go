package game

import "strings"

type LetterStatus string

const (
	Correct LetterStatus = "correct"
	Present LetterStatus = "present"
	Absent  LetterStatus = "absent"
)

type GuessResult struct {
	Word    string         `json:"word"`
	Results []LetterStatus `json:"results"`
}

// compares a guess against the secret word
func CheckGuess(secret, guess string) GuessResult {
	if len(secret) != 5 || len(guess) != 5 {
		return GuessResult{
			Word:    guess,
			Results: make([]LetterStatus, 0),
		}
	}

	secret = strings.ToUpper(secret)
	guess = strings.ToUpper(guess)

	results := make([]LetterStatus, 5)
	secretLetters := []rune(secret)
	guessLetters := []rune(guess)

	used := make([]bool, 5)

	for i := 0; i < 5; i++ {
		if results[i] == Correct {
			continue
		}
		found := false

		for j := 0; j < 5; j++ {
			if !used[j] && guessLetters[i] == secretLetters[j] {
				results[i] = Present
				used[j] = true
				found = true
				break
			}
		}

		if !found {
			results[i] = Absent
		}
	}

	return GuessResult{
		Word:    guess,
		Results: results,
	}
}
