package words

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type WordService struct {
	validWords map[string]bool
	wordList   []string
}

func NewWordService() *WordService {
	ws := &WordService{
		validWords: make(map[string]bool),
		wordList:   make([]string, 0),
	}
	ws.loadWordsFromFile()
	return ws
}

func (ws *WordService) loadWordsFromFile() {
	log.Printf("📚 Loading words from local file")

	file, err := os.Open("internal/words/wordlist.txt")
	if err != nil {
		log.Fatalf("Failed to open word list: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if len(word) == 5 {
			ws.validWords[word] = true
			ws.wordList = append(ws.wordList, word)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading word list: %v", err)
	}

	log.Printf("✅ Loaded %d words", len(ws.wordList))
}

// check if word existss in the list
func (ws *WordService) IsValid(word string) bool {
	return ws.validWords[strings.ToUpper(word)]
}

// get random word
func (ws *WordService) GetRandomWord() string {
	if len(ws.wordList) == 0 {
		return "HELLO"
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(ws.wordList))
	return ws.wordList[randomIndex]
}

// get total number of words
func (ws *WordService) GetWordCount() int {
	return len(ws.wordList)
}
