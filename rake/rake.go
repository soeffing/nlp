package rake

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SplitSentences takes string as input and return slice of sentences (string)
func SplitSentences(text string) []string {
	re, err := regexp.Compile("[.!?,;:]")

	if err != nil {
		fmt.Println(err)
	}
	return re.Split(text, -1)
}

// GenerateCandidateKeywords creates phrases based on the stopword patterns
func GenerateCandidateKeywords(sentences []string, path string) []string {
	re, _ := buildStopWordRegex(path)
	var CandidateWordList []string
	for _, s := range sentences {
		tmp := re.ReplaceAllLiteralString(s, "|")
		phrases := strings.Split(tmp, "|")
		cleanedPhrases := filter(phrases, isNotEmpty)
		for _, p := range cleanedPhrases {
			CandidateWordList = append(CandidateWordList, p)
		}
	}
	return CandidateWordList
}

// CalculateWordScores calculates the scores used to rank the extracted words
// map[string]float64
func CalculateWordScores(phraseList []string) map[string]float64 {
	wordFrecuency := make(map[string]float64)
	wordDegree := make(map[string]float64)
	wordScore := make(map[string]float64)
	for _, phrase := range phraseList {
		// TODO: better word splitting
		wordList := strings.Split(phrase, " ")
		numWords := len(wordList)
		wordListDegree := float64(numWords - 1)
		for _, word := range wordList {
			wordFrecuency[word]++
			wordDegree[word] += wordListDegree
		}

		for k := range wordFrecuency {
			wordDegree[k] = wordDegree[k] + wordFrecuency[k]
		}

		for k := range wordFrecuency {
			wordScore[k] = wordDegree[k] / (wordFrecuency[k] * 1.0)
		}
		//word_score = {}
		//for item in word_frequency:
		//word_score.setdefault(item, 0)
		//word_score[item] = word_degree[item] / (word_frequency[item] * 1.0)  #orig.

	}
	return wordScore
}

// GenerateCandidateKeywordScores generates the final keyword based on scores
func GenerateCandidateKeywordScores(phraseList []string, wordScore map[string]float64) map[string]float64 {
	keywordCandidates := make(map[string]float64)
	for _, phrase := range phraseList {
		wordList := strings.Split(phrase, " ")
		candidateScore := float64(0)
		for _, word := range wordList {
			candidateScore += wordScore[word]
		}
		keywordCandidates[phrase] = candidateScore
	}
	return keywordCandidates
}

// TODO: implement better regex checking
func isNotEmpty(s string) bool {
	if s == "" || s == " " {
		return false
	}
	return true

}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func buildStopWordRegex(path string) (*regexp.Regexp, error) {
	stopwords, _ := loadStopWords(path)

	var regexList []string
	for _, stopword := range stopwords {
		wordRegex := `\b` + stopword + `\W`
		// fmt.Println(wordRegex)
		regexList = append(regexList, wordRegex)
	}

	regexListString := "(?i)" + strings.Join(regexList, "|")
	re, err := regexp.Compile(regexListString)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return re, nil
}

func loadStopWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stopwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stopwords = append(stopwords, scanner.Text())
	}
	return stopwords, scanner.Err()
}

// Run runs the entire rake algorithm and return ranked extracted keywords
func Run(text string) map[string]float64 {
	path, _ := filepath.Abs("rake/data/stopwords.txt")
	sentences := SplitSentences(text)
	phrases := GenerateCandidateKeywords(sentences, path)
	wordScores := CalculateWordScores(phrases)
	candidateKeywords := GenerateCandidateKeywordScores(phrases, wordScores)
	return candidateKeywords

}
