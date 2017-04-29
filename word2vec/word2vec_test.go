package word2vec

import (
	"bufio"
	"fmt"
	"github.com/gonum/matrix/mat64"
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
	// "log"
	"os"
	"strings"
	"testing"
)

var (
	//sentences = [["I", "want", "a", "dog"], ["You", "want", "a", "cat"]]
	sentenceOne   = []string{"I", "want", "a", "dog"}
	sentenceTwo   = []string{"You", "want", "a", "cat"}
	tinySentences = [][]string{sentenceOne, sentenceTwo}
)

func TestBuildVocab(t *testing.T) {
	model := InitModel(1, 1, 10)
	fmt.Println("TestBuildVocab")
	model = BuildVocab(tinySentences, model)
	if len(model.RawVocab) != 6 {
		t.Errorf("Got %d instead of 6 words", len(model.RawVocab))
	}
}

func TestScaleVocab(t *testing.T) {
	fmt.Println("TestScaleVocab")
	model := InitModel(1, 1, 10)
	model = BuildVocab(tinySentences, model)
	model = ScaleVocab(model)
	if len(model.Vocab) != 6 {
		t.Errorf("Got %d instead of 6 words", len(model.Vocab))
	}
}

// func TestFinalizeVocab(t *testing.T) {
// fmt.Println("TestFinalizeVocab")
// model := InitModel(1, 1)
// model = BuildVocab(tinySentences, model)
// model = ScaleVocab(model)
// model = FinalizeVocab(model)
//
// if len(model.CumTable) != 1000 {
// t.Errorf("Got %d instead of 1000 entrries in cumTable", len(model.CumTable))
// }
//
// r, _ := model.Syn1Neg.Dims()
// if r != 6 {
// t.Errorf("Got %d instead of 6 entrries in Syn1Neg", r)
// }
//
// r, _ = model.Syn0.Dims()
// if r != 6 {
// t.Errorf("Got %d instead of 6 entrries in Syn0", r)
// }
// }

// func TestSave(t *testing.T) {
// fmt.Println("TestSave")
// model := InitModel(1, 2)
// model = BuildVocab(tinySentences, model)
// model = ScaleVocab(model)
// model = FinalizeVocab(model)
// model = Train(tinySentences, model)
// model = Save(model, "tiny-model-save.txt")
// }

// func TestTrain(t *testing.T) {
// model := InitModel()
// model = BuildVocab(sentences, model)
// model = ScaleVocab(model)
// model = FinalizeVocab(model)
// Train(sentences, model)
// }

func TestPreMostSimilar(t *testing.T) {
	model := InitModel(1, 5, 50)

	//f, err := os.Open("wiki-data/eng_wikipedia_2012_1M-sentences.txt")
	//f, err := os.Open("wiki-data/eng_news_2015_100K-sentences.txt")
	//f, err := os.Open("wiki-data/eng_wikipedia_2007_10K-sentences.txt")
	f, err := os.Open("wiki-data/4-keyword-corpus.txt")
	//f, err := os.Open("wiki-data/tiny-wiki.txt")

	if err != nil {
		fmt.Printf("error opening model data file: %v\n", err)
		os.Exit(1)
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	// Compiling language specific data into a binary file can be accomplished
	// by using `make <lang>` and then loading the `json` data:
	b, _ := data.Asset("data/english.json")

	// load the training data
	training, _ := sentences.LoadTraining(b)

	// create the default sentence tokenizer
	tokenizer := sentences.NewSentenceTokenizer(training)

	sens := make([][]string, 0)

	for s.Scan() {
		// fmt.Println(s.Text())
		// fmt.Println(strings.Fields(s.Text()))
		sentences := tokenizer.Tokenize(s.Text())
		for _, s := range sentences {
			// fmt.Println(strings.Fields(s.Text))
			sens = append(sens, strings.Fields(s.Text))
		}
	}

	fmt.Println("Train on %d sentences", len(sens))
	model = BuildVocab(sens, model)
	model = ScaleVocab(model)
	model = FinalizeVocab(model)

	terms := make([]string, 0)
	// terms = append(terms, "vistas")
	// terms = append(terms, "configuration")
	// terms = append(terms, "cuisine")
	// terms = append(terms, "macbook")
	terms = append(terms, "computer")

	for _, term := range terms {
		ans := MostSimilar(term, 20, model)
		fmt.Printf("Most Similar for %s\n", term)
		fmt.Println(ans)
	}
}

func TestTrainCorpus(t *testing.T) {
	fmt.Println("TestTrainCorpus")
	model := InitModel(5, 20, 50)

	//f, err := os.Open("wiki-data/eng_wikipedia_2012_1M-sentences.txt")
	//f, err := os.Open("wiki-data/eng_news_2015_100K-sentences.txt")
	//f, err := os.Open("wiki-data/eng_wikipedia_2007_10K-sentences.txt")
	//f, err := os.Open("wiki-data/4-keyword-corpus.txt")
	f, err := os.Open("wiki-data/tiny-wiki.txt")

	if err != nil {
		fmt.Printf("error opening model data file: %v\n", err)
		os.Exit(1)
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	// Compiling language specific data into a binary file can be accomplished
	// by using `make <lang>` and then loading the `json` data:
	b, _ := data.Asset("data/english.json")

	// load the training data
	training, _ := sentences.LoadTraining(b)

	// create the default sentence tokenizer
	tokenizer := sentences.NewSentenceTokenizer(training)

	sens := make([][]string, 0)

	for s.Scan() {
		// fmt.Println(s.Text())
		// fmt.Println(strings.Fields(s.Text()))
		sentences := tokenizer.Tokenize(s.Text())
		for _, s := range sentences {
			// fmt.Println(strings.Fields(s.Text))
			sens = append(sens, strings.Fields(s.Text))
		}
	}

	fmt.Println("Train on %d sentences", len(sens))
	model = BuildVocab(sens, model)
	model = ScaleVocab(model)
	model = FinalizeVocab(model)
	model = Train(sens, model)

	//model = Save(model, "1M_model_final.txt")
	//model = Save(model, "100k_model_5-min-count_5_epochs.txt")
	//model = Save(model, "big_50_epochs.txt")
	//model = Save(model, "4-keyword_model.txt")
	model = Save(model, "tiny_model.txt")

}

// func TestMostSimilar(t *testing.T) {
// fmt.Println("TestMostSimilar")
// model := Load("big_50_epochs.txt")
// fmt.Println("Vocab size: %i", len(model.Vocab))
//
// MostSimilar("can", model)
//
// }

// func TestLoad(t *testing.T) {
// model := Load("1M_model.txt")
// //model := Load("100k_model.txt")
//
// fmt.Println("Vocab size: %i", len(model.Vocab))
// // if len(model.RawVocab) != 6 {
// // t.Errorf("Got %d instead of 6 words", len(model.RawVocab))
// // }
//
// sim0 := Similarity("she", "woman", model)
// fmt.Println("sim she vs woman")
// fmt.Println(sim0)
//
// sim1 := Similarity("man", "he", model)
// fmt.Println("sim man vs he")
// fmt.Println(sim1)
//
// sim1 = Similarity("man", "man", model)
// fmt.Println("sim man vs man")
// fmt.Println(sim1)
//
// }

func TestMostSimilar(t *testing.T) {
	//model := Load("1M_model.txt")
	model := Load("1M_model_final.txt")
	//model := Load("gensim-model.txt")
	//model := Load("4-keyword_model.txt")

	fmt.Println("Vocab size: %i", len(model.Vocab))

	// sim0 := Similarity("she", "woman", model)
	// fmt.Println("sim she vs woman")
	// fmt.Println(sim0)
	//
	// sim1 := Similarity("man", "he", model)
	// fmt.Println("sim man vs he")
	// fmt.Println(sim1)
	//
	// sim1 = Similarity("man", "man", model)
	// fmt.Println("sim man vs man")
	// fmt.Println(sim1)

	terms := make([]string, 0)
	terms = append(terms, "queen")
	terms = append(terms, "man")
	terms = append(terms, "woman")
	terms = append(terms, "king")

	for _, term := range terms {
		ans := MostSimilar(term, 100, model)
		fmt.Printf("Most Similar for %s\n", term)
		fmt.Println(ans)
	}
}

func TestMostSimilarByVector(t *testing.T) {
	model := Load("1M_model_final.txt")
	//model := Load("gensim-model.txt")
	//model := Load("4-keyword_model.txt")

	terms := make([]string, 0)
	terms = append(terms, "queen")
	terms = append(terms, "man")
	terms = append(terms, "woman")
	terms = append(terms, "king")

	for _, term := range terms {
		id := model.Word2Index[term]
		phrase := model.Vocab[id]
		fmt.Printf("Word: %v\n", phrase.Literal)
		ans := MostSimilarByVector(phrase.Vector, 20, model)
		fmt.Println(ans)
	}

}

func TestVectorCalcs(t *testing.T) {
	fmt.Println("TestVectorCalcs")
	model := Load("100k_model_5-min-count_5_epochs.txt")

	// queen ≈ king − man + woman
	// washington ≈ berlin − germany + usa

	kingVector := model.Vocab[model.Word2Index["king"]].Vector
	manVector := model.Vocab[model.Word2Index["man"]].Vector
	womanVector := model.Vocab[model.Word2Index["woman"]].Vector

	manWomanVec := mat64.NewVector(50, nil)
	manWomanVec.AddVec(manVector, womanVector)

	finalVec := mat64.NewVector(50, nil)
	finalVec.SubVec(kingVector, manWomanVec)

	ans := MostSimilarByVector(finalVec, 50, model)
	fmt.Println(ans)

	//model := Load("gensim-model.txt")
	//model := Load("4-keyword_model.txt")
}

// func TestTrainBig(t *testing.T) {
// model := InitModel()
//
// file, err := os.Open("small_file.txt")
// if err != nil {
// log.Fatal(err)
// }
// defer file.Close()
//
// // var bigSentences [][]string
// bigSentences := make([][]string, 0)
//
// scanner := bufio.NewScanner(file)
// for scanner.Scan() {
// bigSentences = append(bigSentences, strings.Fields(scanner.Text()))
// // fmt.Println(strings.Fields(scanner.Text()))
// }
//
// model = BuildVocab(bigSentences, model)
//
// model = ScaleVocab(model)
// model = FinalizeVocab(model)
//
// model = Train(bigSentences, model)
//
// if err := scanner.Err(); err != nil {
// log.Fatal(err)
// }
//
// sim := Similarity("the", "of", model)
// fmt.Println("sim the vs of")
// fmt.Println(sim)
//
// sim2 := Similarity("egyptian", "egypt", model)
// fmt.Println("sim egyptian vs egypt")
// fmt.Println(sim2)
//
// sim3 := Similarity("egypt", "egyptian", model)
// fmt.Println("sim egypt vs egyptian")
// fmt.Println(sim3)
//
// sim4 := Similarity("prehistoric", "egypt", model)
// fmt.Println("sim prehistoric vs egypt")
// fmt.Println(sim4)
//
// sim5 := Similarity("see", "egypt", model)
// fmt.Println("sim see vs egypt")
// fmt.Println(sim5)
//
// sim6 := Similarity("to", "of", model)
// fmt.Println("sim to vs of")
// fmt.Println(sim6)
//
// sim7 := Similarity("to", "the", model)
// fmt.Println("sim to vs the")
// fmt.Println(sim7)
//
// sim8 := Similarity("government", "dynasty", model)
// fmt.Println("sim government vs dynasty")
// fmt.Println(sim8)
//
// sim9 := Similarity("government", "tourist", model)
// fmt.Println("sim government vs tourist")
// fmt.Println(sim9)
//
// sim10 := Similarity("government", "military", model)
// fmt.Println("sim government vs military")
// fmt.Println(sim10)
//
// loadedModel := Load("dat2.txt")
//
// if len(loadedModel.RawVocab) != 6 {
// t.Errorf("Got %d instead of 6 words for the loaded model", len(model.RawVocab))
// }
//
// }
