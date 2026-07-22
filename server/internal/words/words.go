package words

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	WordAPIURL = "https://wordle-api.cyclic.app/words"
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
	ws.loadWordFromAPI()
	return ws
}

func (ws *WordService) loadWordFromAPI() {
	resp, err := http.Get(WordAPIURL)
	if err != nil {
		log.Fatal("Failed to fetch word list from API: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read API response: %v", err)
	}

	var words []string

	if err := json.Unmarshal(body, &words); err != nil {
		log.Fatalf("Failed to parse word list: %v", err)
	}

	for _, word := range words {
		word = strings.ToUpper(strings.Trim(word))
		if len(word) == 5 {
			ws.validWords[word] = true
			ws.wordList = append(ws.wordList, word)
		}
	}
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
