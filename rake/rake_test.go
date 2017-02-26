package rake

import (
	"testing"
)

var (
	text      = "Compatibility of systems of linear constraints over the set of natural numbers. Criteria of compatibility of a system of linear Diophantine equations, strict inequations, and nonstrict inequations are considered. Upper bounds for components of a minimal set of solutions and algorithms of construction of minimal generating sets of solutions for all types of systems are given. These criteria and the corresponding algorithms for constructing a minimal supporting set of solutions can be used in solving all the considered types of systems and systems of mixed types."
	smallText = "This is a tiny test. A test for rake, a keyword extraction technology."
)

func TestSplitSentences(t *testing.T) {
	expectedLength := 7
	sentences := SplitSentences(text)

	if len(sentences) != expectedLength {
		t.Fatalf("SplitSentences does not split the sentences correctly")
	}
}

func TestGenerateCandidateKeywords(t *testing.T) {
	sentencesList := SplitSentences(smallText)
	actual := GenerateCandidateKeywords(sentencesList, "data/stopwords.txt")

	if actual[3] != "keyword extraction technology" {
		t.Fatalf("Candidate keywords are not being generated correctly")
	}
}

func TestCalculateWordScores(t *testing.T) {
	var phraseList []string
	phraseList = append(phraseList, "rake")
	phraseList = append(phraseList, "keyword extraction technology")

	actual := CalculateWordScores(phraseList)

	if actual["extraction"] != 3 && actual["rake"] != 2 {
		t.Fatalf("CalculateWordScores not working as expected")
	}

}

func TestGenerateCandidateKeywordScores(t *testing.T) {
	var phraseList []string
	phraseList = append(phraseList, "rake")
	phraseList = append(phraseList, "keyword extraction technology")

	wordScore := make(map[string]float64)
	wordScore["rake"] = 2
	wordScore["keyword"] = 3
	wordScore["extraction"] = 3
	wordScore["technology"] = 3

	actual := GenerateCandidateKeywordScores(phraseList, wordScore)

	if actual["rake"] != 2 && actual["keyword extraction technology"] != 9 {
		t.Fatalf("GenerateCandidateKeywordScores not working as expected")
	}
}
