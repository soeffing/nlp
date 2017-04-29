package word2vec

import (
	"fmt"
	// "github.com/a-h/round"
	// "github.com/soeffing/nlp/matrix"
	// "github.com/soeffing/nlp/vector"
	// "github.com/gonum/blas"
	// "github.com/gonum/blas/blas64"
	// "github.com/gonum/lapack/lapack64"
	// "github.com/gonum/matrix"
	"bufio"
	"sort"
	// "bytes"
	// "encoding/binary"
	//"github.com/a-h/round"
	"github.com/gonum/matrix/mat64"
	// "github.com/gonum/stat"
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	//	"github.com/gonum/matrix"
)

// Trains a word2vec model after Mikolov: https://arxiv.org/abs/1310.4546
// https://papers.nips.cc/paper/5021-distributed-representations-of-words-and-phrases-and-their-compositionality.pdf

// python implementation: https://github.com/RaRe-Technologies/gensim/blob/develop/gensim/models/word2vec.py

// Model contains the word2vec vocab, vectors, parameters, etc.
type Model struct {
	RawVocab         map[string]*Phrase
	Vocab            []*Phrase
	Word2Index       map[string]int
	MinCount         int
	Sample           float64
	CumTable         []int
	Syn0             *mat64.Dense
	Syn1Neg          *mat64.Dense // size: vocab_size * vector_dimension
	Syn0Norm         *mat64.Dense
	Window           int
	Negative         int
	TotalCorpusCount int
	Alpha            float64
	Epochs           int
	VecDim           int
	seed             int64
	stopwords        []string
}

// Phrase contains the literal text, vecotr and count of a given phrase
type Phrase struct {
	Literal string
	Vector  *mat64.Vector
	Count   int
	Id      int
	// used for subsampling
	SampleInt int
	// used for populating NegativeSamplingTable
	Probability float64
	Updated     int64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Vocab []*Phrase

func (p Vocab) Len() int           { return len(p) }
func (p Vocab) Less(i, j int) bool { return p[i].Updated < p[j].Updated }
func (p Vocab) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Gensim run dowm

// func train_sg_pair()

// InitModel initiliazes the model
func InitModel(epochs int, minCount int, dim int) *Model {
	stopwords := make([]string, 0)

	f, err := os.Open("stopwords.txt")
	if err != nil {
		os.Exit(1)
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		stopwords = append(stopwords, s.Text())
	}

	// fmt.Printf("stopwords: %v\n", stopwords)

	model := &Model{
		RawVocab:         make(map[string]*Phrase),
		Vocab:            make(Vocab, 0),
		MinCount:         minCount,
		Sample:           0.001,
		CumTable:         make([]int, 0),
		Word2Index:       make(map[string]int),
		Window:           7,
		Negative:         5,
		TotalCorpusCount: 0,
		Alpha:            0.025,
		Epochs:           epochs,
		VecDim:           dim,
		seed:             7456393,
		stopwords:        stopwords,
	}
	return model
}

// Similarity returns the cosine similarity between two words
// def cosine_measure(v1, v2):
// prod = dot_product(v1, v2)
// len1 = math.sqrt(dot_product(v1, v1))
// len2 = math.sqrt(dot_product(v2, v2))
// return prod / (len1 * len2)
func Similarity(vector1 *mat64.Vector, vector2 *mat64.Vector, model *Model) float64 {
	// fmt.Println(w1)
	//index1 := model.Word2Index[w1]
	//index2 := model.Word2Index[w2]

	//fmt.Printf("index1: %d\n", index1)
	//fmt.Printf("index2: %d\n", index2)

	// // fmt.Println(model.Syn0Norm.Dims())

	// vector1 := model.Vocab[index1].Vector
	// vector2 := model.Vocab[index2].Vector
	//fmt.Printf("vector1: %v\n", vector1)
	//fmt.Printf("vector2: %v\n", vector2)

	// neu1e := mat64.NewVector(300, nil)

	dotProduct := mat64.Dot(vector1, vector2)

	//fmt.Println("dotProduct: %v\n", dotProduct)

	len1 := math.Sqrt(mat64.Dot(vector1, vector1))
	len2 := math.Sqrt(mat64.Dot(vector2, vector2))

	//fmt.Println("len1: %d\n", len1)
	//fmt.Println("len2: %d\n", len2)

	res := dotProduct / (len1 * len2)

	// fmt.Println("res: %d\n", res)

	return res
}

type SimPair struct {
	Key string
	Sim float64
}

type SimPairList []SimPair

func (p SimPairList) Len() int           { return len(p) }
func (p SimPairList) Less(i, j int) bool { return p[i].Sim < p[j].Sim }
func (p SimPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func MostSimilar(positive string, top int, model *Model) []SimPair {
	// fmt.Println(positive)
	spl := make(SimPairList, len(model.Vocab))
	fmt.Printf("Positive: %v\n", positive)
	for idx, phraseObj := range model.Vocab {
		if phraseObj.Literal != positive {
			//fmt.Printf("phraseObj.Literal: %v\n", phraseObj.Literal)

			index1 := model.Word2Index[positive]
			index2 := model.Word2Index[phraseObj.Literal]

			vector1 := model.Vocab[index1].Vector
			vector2 := model.Vocab[index2].Vector

			sim := Similarity(vector1, vector2, model)
			spl[idx] = SimPair{phraseObj.Literal, sim}
			// fmt.Println("%v has %v", phraseObj.Literal, sim)
		}
	}

	sort.Sort(sort.Reverse(spl))

	res := make(SimPairList, top+1)
	for c, pair := range spl {
		res[c] = SimPair{pair.Key, pair.Sim}
		if c >= top {
			break
		}
	}
	return res
}

func MostSimilarByVector(vec *mat64.Vector, top int, model *Model) []SimPair {
	// fmt.Println(positive)
	spl := make(SimPairList, len(model.Vocab))
	// fmt.Printf("vec: %v\n", vec)
	for idx, phraseObj := range model.Vocab {
		//if phraseObj.Literal != positive {
		//fmt.Printf("phraseObj.Literal: %v\n", phraseObj.Literal)
		index2 := model.Word2Index[phraseObj.Literal]
		vector2 := model.Vocab[index2].Vector

		sim := Similarity(vec, vector2, model)
		spl[idx] = SimPair{phraseObj.Literal, sim}
		// fmt.Println("%v has %v", phraseObj.Literal, sim)
		//}
	}

	sort.Sort(sort.Reverse(spl))

	res := make(SimPairList, top+1)
	for c, pair := range spl {
		res[c] = SimPair{pair.Key, pair.Sim}
		if c >= top {
			break
		}
	}
	return res
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// BuildVocab cycles through sentences and builds vocab
func BuildVocab(sentences [][]string, model *Model) *Model {
	for _, sentence := range sentences {
		for _, word := range sentence {

			temp := cleanString(word)
			cleanWord := strings.Replace(temp, " ", "", -1)
			cleanWord = strings.ToLower(cleanWord)

			//if stringInSlice(cleanWord, model.stopwords) == true {
			//continue
			//}

			if phraseObj, ok := model.RawVocab[cleanWord]; ok {

				phraseObj.Count++
			} else {
				phrase := &Phrase{
					Literal: cleanWord,
					Count:   1,
					Updated: 0,
				}
				// // fmt.Println(model.Vocab[word])
				model.RawVocab[cleanWord] = phrase
			}
			// //fmt.Println(word)
		}
	}

	return model
}

// ScaleVocab applies MinCount
func ScaleVocab(model *Model) *Model {
	// retainTotal := 0
	idx := 0
	for word, phraseObj := range model.RawVocab {
		if phraseObj.Count >= model.MinCount {
			phraseObj.Id = idx
			model.Vocab = append(model.Vocab, phraseObj)
			model.Word2Index[word] = idx
			idx++
			// model.Vocab[phrase] = phraseObj
			// retainTotal += phraseObj.Count

			// incease total corpus count
			model.TotalCorpusCount += phraseObj.Count
		}
	}

	// TODO: Double check those calculation with http://mccormickml.com/2017/01/11/word2vec-tutorial-part-2-negative-sampling/

	thresholdCount := model.Sample * float64(model.TotalCorpusCount)

	downsampleTotal := 0.0
	downsampleUnique := 0.0

	for index, phraseObj := range model.Vocab {
		v := float64(model.RawVocab[phraseObj.Literal].Count)
		wordProbability := (math.Sqrt(v/thresholdCount) + 1) * (thresholdCount / v)

		if wordProbability < 1.0 {
			downsampleUnique++
			downsampleTotal += wordProbability * v
		} else {
			wordProbability = 1.0
			downsampleTotal += v
		}
		// // fmt.Println(round.ToEven(float64(wordProbability*math.Pow(2, 32)), 20))
		// // fmt.Println(int(round.ToEven(float64(wordProbability*math.Pow(2, 32)), 20)))
		// // fmt.Println("SampleInt")
		// // fmt.Println(wordProbability * math.Pow(2, 32))
		// // fmt.Println(int(wordProbability * math.Pow(2, 32)))
		// // fmt.Println(round.ToEven(float64(wordProbability*math.Pow(2, 32)), 20))
		model.Vocab[index].SampleInt = int(wordProbability * math.Pow(2, 32))
		// fmt.Printf("word: %s\n", phraseObj.Literal)
		// fmt.Printf("count: %d\n", phraseObj.Count)
		// fmt.Printf("Sample int: %d\n", model.Vocab[index].SampleInt)
		// fmt.Printf("Probability: %f\n", wordProbability)
		// fmt.Printf("#####\n")

	}

	return model
}

// FinalizeVocab creates the weigths and the
func FinalizeVocab(model *Model) *Model {
	// sort vocab?
	// model = makeCumTable(model)
	model = makeCumTableAlternative(model)
	// fmt.Printf("model.CumTable: %v\n", model.CumTable)
	model = resetWeights(model)

	return model
}

// http://mccormickml.com/2017/01/11/word2vec-tutorial-part-2-negative-sampling/
// The way this selection is implemented in the C code is interesting.
// They have a large array with 100M elements (which they refer to as the unigram table).
// They fill this table with the index of each word in the vocabulary multiple times,
// and the number of times a word’s index appears in the table is given by P(wi) * table_size.
// Then, to actually select a negative sample, you just generate a random integer between 0 and 100M,
// and use the word at that index in the table. Since the higher probability words occur more times in the table,
// you’re more likely to pick those.

// Equation: probability of word = word count of word ^0.75 / word count of entire corpus ^0.75
// Equation 2: probability of word * table_size
// Table will be an array of 1000 fields (for dev environment) and 100 million fields production environment

func makeCumTableAlternative(model *Model) *Model {

	// fmt.Println("inside makeCumTableAlternative")
	trainWordsPow := 0.0
	power := 0.75

	for index := range model.Vocab {
		trainWordsPow += math.Pow(float64(model.Vocab[index].Count), 0.75)
	}

	i := 0
	fmt.Printf("size model.Vocab %v\n", len(model.Vocab))
	d1 := math.Pow(float64(model.Vocab[i].Count), power) / trainWordsPow

	fmt.Printf("d1: %f\n", d1)

	tableSize := 100000000

	// TODO: set proper table size
	for a := 1; a <= tableSize; a++ {
		// model.CumTable = append(model.CumTable, i)
		// [Chris] - The table may contain multiple elements which hold value 'i'.
		//
		model.CumTable = append(model.CumTable, i)
		// table[a] = i;
		// [Chris] - If the fraction of the table we have filled is greater than the
		//           fraction of this words weight / all word weights, then move to the next word.
		// fmt.Printf("float64(a/tableSize): %f\n", float64(a)/float64(tableSize))
		if float64(a)/float64(tableSize) > d1 {
			// fmt.Printf("increase i at index: %d\n", a)
			if i == len(model.Vocab)-1 {
				// fmt.Println("lower i")
				i = len(model.Vocab) - 1
			} else {
				i++
			}

			// fmt.Printf("i: %d\n", i)
			// fmt.Printf("vocab size: %d\n", len(model.Vocab))

			d1 += math.Pow(float64(model.Vocab[i].Count), power) / trainWordsPow
		}

	}

	return model
}

func resetWeights(model *Model) *Model {
	// Reset all projection weights to an initial (untrained) state, but keep the existing vocabulary

	// model.Syn0 = mat64.NewDense(300, len(model.Vocab), data) // matrix.NewMatrix(300, len(model.Vocab))
	// model.Syn1Neg = matrix.NewMatrix(300, len(model.Vocab))

	// d := make([]float64, n*n)

	// vecs := make([]vector.Vector, 0)

	model.Syn0 = mat64.NewDense(len(model.Vocab), model.VecDim, nil)

	for index, phraseObj := range model.Vocab {
		randomVector := seededVector(model, model.Vocab[index].Literal)
		// // fmt.Println("randomVector")
		// // fmt.Println(randomVector)
		// // fmt.Println("index")
		// // fmt.Println(index)
		//
		// //fmt.Println("model.Syn0")
		// //fmt.Println(len(model.Syn0))

		// mat64.Dense.

		// data[index] := randomVector

		model.Syn0.SetRow(index, randomVector)

		phraseObj.Vector = model.Syn0.RowView(index)
		// model.
		// model.Syn0.InitRow(index, randomVector)

		// create Zero vector and add it to the Syn1Neg
		// zeroElements := make([]float64, 0)

		// for i := 1; i <= 300; i++ {
		// zeroElements = append(zeroElements, 0.0)
		// }

		// zeroVector := mat64.NewVector(300, zeroElements) // NewVector(300, zeroElements...)

	}

	model.Syn1Neg = mat64.NewDense(len(model.Vocab), model.VecDim, nil)

	// layer1_size = vector dimension
	// zeros(columns, rows)
	// syn1neg = vocab_size * vector_dimension
	// self.syn1neg = zeros((len(self.wv.vocab), self.layer1_size), dtype=REAL)

	// // fmt.Println("vecs")
	// // fmt.Println(vecs)

	// model.Syn0.Init(vecs)

	return model
	// randomize weights vector by vector, rather than materializing a huge random matrix in RAM at once
	// for i in xrange(len(self.wv.vocab)):
	// # construct deterministic seed from word AND seed argument
	// self.wv.syn0[i] = self.seeded_vector(self.wv.index2word[i] + str(self.seed))

	// self.wv.syn0 = empty((len(self.wv.vocab), self.vector_size), dtype=REAL)
	// # randomize weights vector by vector, rather than materializing a huge random matrix in RAM at once
	// for i in xrange(len(self.wv.vocab)):
	// # construct deterministic seed from word AND seed argument
	// self.wv.syn0[i] = self.seeded_vector(self.wv.index2word[i] + str(self.seed))
	// if self.hs:
	// self.syn1 = zeros((len(self.wv.vocab), self.layer1_size), dtype=REAL)
	// if self.negative:
	// self.syn1neg = zeros((len(self.wv.vocab), self.layer1_size), dtype=REAL)
	// self.wv.syn0norm = None
	//
	// self.syn0_lockf = ones(len(self.wv.vocab), dtype=REAL)  # zeros suppress learning

}

// func round(num float64) int {
// return int(num + math.Copysign(0.5, num))
// }

// func toFixed(num float64, precision int) float64 {
// output := math.Pow(10, float64(precision))
// return float64(round(num*output)) / output
// }

func seededVector(model *Model, seedString string) []float64 {
	// vector.Vector = vector.NewVector(300)
	// vector = vector.NewVector(300)

	// oneRand := hash(seedString)
	// // fmt.Println(oneRand)
	vectorElements := make([]float64, 0)
	//vectorSlice := make([]float64, 0)
	seed := model.seed + rand.Int63()
	// fmt.Printf("seed: %v\n", seed)
	rand.Seed(seed)

	for i := 1; i <= model.VecDim; i++ {
		//vectorElements = append(vectorElements, (rand.Float64()-0.5)/float64(model.VecDim))
		vectorElements = append(vectorElements, rand.Float64()-0.5)
	}

	// fmt.Println("vectorElements")
	// fmt.Println(vectorElements)

	// actualVector := vector.NewVector(300, vectorElements...)
	// actualVector.Zero()
	// actualVector := mat64.NewVector(300, vectorElements)
	//actualVector := new(vector.Vector)

	// actualVector.InitV(vec)

	// // fmt.Println("actualVector after NewVector")
	// // fmt.Println(actualVector)

	// // fmt.Println("vectorElements")
	// // fmt.Println(vectorElements)

	//newVec = actualVector.Init(vectorElements...)

	// actualVector.Set(1, 1.5)

	// // fmt.Println("actualVector")
	// // fmt.Println(actualVector)
	// // fmt.Println(newVec)
	// //fmt.Println(hash(seedString))

	// Create one 'random' vector (but deterministic by seed_string)
	// Note: built-in hash() may vary by Python version or even (in Py3.x) per launch
	// once = random.RandomState(self.hashfxn(seed_string) & 0xffffffff)
	// return (once.rand(self.vector_size) - 0.5) / self.vector_size
	return vectorElements

}

// From: http://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Save(model *Model, path string) *Model {
	f, err := os.Create(path)
	check(err)

	defer f.Close()

	// modelData := make([]int, 0)

	// fmt.Println("len(model.Vocab)")
	// fmt.Println(len(model.Vocab))

	stringDim := strconv.Itoa(model.VecDim)
	stringSize := strconv.Itoa(len(model.Vocab))

	f.WriteString(stringDim + " " + stringSize)

	f.WriteString("\n")

	// newReader := bufio.NewReader(f)
	newWriter := bufio.NewWriter(f)

	for _, voc := range model.Vocab {
		vec := model.Syn0.RowView(voc.Id)
		newWriter.WriteString(voc.Literal + "\n")

		for _, elem := range vec.RawVector().Data {
			newWriter.WriteString(strconv.FormatFloat(elem, 'g', 15, 64) + " ")
			// err := binary.Write(buf, binary.BigEndian, elem)
			// if err != nil {
			// // fmt.Println("binary.Write failed:", err)
			// }
		}
		newWriter.WriteString("\n")
		// f.Write(buf.Bytes())

	}

	newWriter.Flush()

	// write debug data: update counts of words
	f, err = os.Create("debug/sorted-updates.txt")
	check(err)

	defer f.Close()

	newWriter = bufio.NewWriter(f)

	vocab := make(Vocab, len(model.Vocab))
	for idx, phraseObj := range model.Vocab {
		phrase := &Phrase{
			Literal: phraseObj.Literal,
			Count:   phraseObj.Count,
			Updated: phraseObj.Updated,
		}
		vocab[idx] = phrase
	}

	sort.Sort(sort.Reverse(vocab))

	for _, voc := range vocab {
		// vec := model.Syn0Norm.RowView(voc.Id)
		newWriter.WriteString(voc.Literal + " " + strconv.FormatInt(voc.Updated, 10) + " " + strconv.FormatInt(int64(voc.Count), 10))
		newWriter.WriteString("\n")
		// f.Write(buf.Bytes())
	}

	newWriter.Flush()

	//modelData = append(modelData, 300)
	//modelData = append(modelData, len(model.Vocab))
	//
	// //fmt.Println("modelData")
	// //fmt.Println(modelData)

	// buf := new(bytes.Buffer)

	//for _, elem := range modelData {
	// //fmt.Println("elem")
	// //fmt.Println(elem)
	//
	//r := bytes.NewReader([]byte{elem})
	//f.Write()
	//
	//
	//// binary.Write(buf, binary.BigEndian, elem)
	//
	//}

	// for _, voc := range model.Vocab {
	// vec := model.Syn0Norm.RowView(voc.Id)
	//
	// for _, elem := range vec.RawVector().Data {
	// f.WriteString(voc.Literal)
	// err := binary.Write(buf, binary.BigEndian, elem)
	// if err != nil {
	// // fmt.Println("binary.Write failed:", err)
	// }
	// }
	// f.Write(buf.Bytes())
	// f.WriteString("/n")
	//
	// }
	return model
}

func Load(modelPath string) *Model {

	model := InitModel(50, 5, 0)

	f, err := os.Open(modelPath)
	if err != nil {
		// fmt.Printf("error opening model data file: %v\n", err)
		os.Exit(1)
	}

	// br := bufio.NewReader(f)

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	size_dim := make([]string, 0)
	words := make([]string, 0)
	wordVectorSlice := make([]*Phrase, 0)
	counter := 0
	for s.Scan() {
		if counter == 0 {
			size_dim = strings.Fields(s.Text())
			fmt.Printf("Size dim: %v\n", size_dim[0])
			intDim, _ := strconv.Atoi(size_dim[0])
			model.VecDim = intDim

		} else if counter%2 == 0 {
			// // fmt.Println("Even number %i", counter)
			strVector := strings.Fields(s.Text())
			floatVector := make([]float64, 0)
			for _, strFloat := range strVector {
				if n, err := strconv.ParseFloat(strFloat, 64); err == nil {
					floatVector = append(floatVector, n)
				}
			}

			vector := mat64.NewVector(model.VecDim, floatVector)
			// // fmt.Println("word before setting map: %v", words[len(words)-1])

			phrase := &Phrase{
				Literal: words[len(words)-1],
				Count:   0,
				Vector:  vector,
			}

			wordVectorSlice = append(wordVectorSlice, phrase)

		} else {
			// // fmt.Println("Uneven number %i", counter)
			word := s.Text()
			words = append(words, word)
			// // fmt.Println(word)
			// wordVectorMap[word] :=
		}

		// // fmt.Println(s.Text())
		// // fmt.Println("----new line ----")
		counter++
	}
	// // fmt.Println("wordVectorMap: %v", wordVectorMap)

	for index, phraseObj := range wordVectorSlice {
		model.Vocab = append(model.Vocab, phraseObj)
		model.Word2Index[phraseObj.Literal] = index
	}

	// index := model.Word2Index["BC,"]
	// // fmt.Println("wordVectorMap: %v", model.Vocab[index].Vector)

	// // fmt.Println(model.Vocab[word])
	// model.RawVocab[word] = phrase

	// if err != nil {
	// // fmt.Printf("error opening model data file: %v\n", err)
	// os.Exit(1)
	// }
	//
	// defer f.Close()
	//
	// var size, dim int
	// _, errDim := fmt.Fscanf(f, "%d", &dim)
	// _, errSize := fmt.Fscanf(f, "%d", &size)
	//
	// // fmt.Println("dim")
	// // fmt.Println(dim)
	//
	// if errDim != nil {
	// // fmt.Println("error reading file for dim")
	// // fmt.Println(errDim.Error())
	// return nil
	// }
	//
	// if errSize != nil {
	// // fmt.Println("error reading file for size")
	// // fmt.Println(errSize.Error())
	// return nil
	// }

	// newReader := bufio.NewReader(f)
	// newWriter := bufio.NewWriter(f)
	// newWriter.Write(1.999999999)

	// newReader.ReadString('\n')

	// var word string
	// _, errWord := fmt.Fscanf(f, "%d", &word)
	//
	// if errWord != nil {
	// // fmt.Println("error reading first word")
	// // fmt.Println(errWord.Error())
	// return nil
	// }

	// if n != 2 {
	// // fmt.Println("could not extract size/dim from binary model data")
	// return nil
	// }

	// // fmt.Println("dim")
	// // fmt.Println(dim)
	// // fmt.Println("size")
	// // fmt.Println(size)

	// vocab := make([][]float64, 0)

	//v := make([]float64, 300)
	//for i := 1; i <= 26; i++ {
	//if err := binary.Read(br, binary.BigEndian, v); err != nil {
	// //fmt.Println(err.Error())
	//return nil
	//}
	//
	// //fmt.Println(v)
	//vocab = append(vocab, v)
	//}
	//
	// //fmt.Println("Actual vector count")
	// //fmt.Println(len(vocab))
	//

	// m := &Model{
	// words: make(map[string]Vector, size),
	// dim:   dim,
	// }

	return model
}

const delim = "?!.;,*\"()'"

func isDelim(c string) bool {
	if strings.Contains(delim, c) {
		return true
	}
	return false
}

func cleanString(input string) string {

	size := len(input)
	temp := ""
	var prevChar string

	for i := 0; i < size; i++ {
		// //fmt.Println(input[i])
		str := string(input[i]) // convert to string for easier operation
		if (str == " " && prevChar != " ") || !isDelim(str) {
			temp += str
			prevChar = str
		} else if prevChar != " " && isDelim(str) {
			temp += " "
		}
	}
	return temp
}

// Train train the word2vec model over the corpus
func Train(sentences [][]string, model *Model) *Model {

	negLabels := make([]float64, model.Negative+1)
	negLabels = append(negLabels, 1.0)
	for len(negLabels) < model.Negative+1 {
		negLabels = append(negLabels, 0.0)
	}

	startingAlpha := model.Alpha
	skippedWords := make([]string, 0)

	rand.Seed(model.seed)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)

	wordCount := 1

	for epoch := 1; epoch <= model.Epochs; epoch++ {
		for _, sentence := range sentences {
			if len(sentence) == 0 {
				continue
			}
			trainSentence := make([]*Phrase, 0)

			for _, w := range sentence {

				temp := cleanString(w)
				word := strings.Replace(temp, " ", "", -1)
				word = strings.ToLower(word)
				if idx, ok := model.Word2Index[word]; ok {
					wordCount++

					ran := (math.Sqrt(float64(model.Vocab[idx].Count)/(model.Sample*float64(wordCount))) + 1) * (model.Sample * float64(wordCount)) / float64(model.Vocab[idx].Count)
					nextRandom := r2.Float64()

					if ran > nextRandom {
						trainSentence = append(trainSentence, model.Vocab[idx])
					} else {
						skippedWords = append(skippedWords, word)
					}
				} else {
					continue
				}

			}

			sentencePosition := 0

			if len(trainSentence) == 0 {
				continue
			}

			for sentencePosition < len(trainSentence) {
				residual := wordCount % 10000
				if residual < 10 {

					model.Alpha = startingAlpha * (1.0 - float64(wordCount)/(float64(epoch)*float64(model.TotalCorpusCount)+1.0))

					if model.Alpha < startingAlpha*0.0001 {
						model.Alpha = startingAlpha * 0.0001
					}
					fmt.Printf("new Alpha: %v\n", model.Alpha)
					model.Alpha = 0.025
				}

				word := trainSentence[sentencePosition]
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)

				nextRandom := r1.Intn(10000)
				b := int(math.Mod(float64(nextRandom), float64(model.Window)))

				for a := b; a < model.Window*2+1-b; a++ {
					c := sentencePosition - model.Window + a

					if c < 0 {
						continue
					}

					if c >= len(trainSentence) {
						continue
					}

					lastWord := trainSentence[c]

					if lastWord.Literal == word.Literal {
						continue
					}

					neu1e := mat64.NewVector(model.VecDim, nil)

					if word.Id == lastWord.Id {
						continue
					}

					l1 := model.Syn0.RowView(lastWord.Id)

					target := 0
					label := 0
					for d := 0; d < model.Negative+1; d++ {

						if d == 0 {
							target = word.Id
							label = 1
						} else {

							random := int(r1.Intn(len(model.CumTable)))

							target = model.CumTable[random]

							if target == word.Id {
								continue
							}

							label = 0
						}

						l2 := model.Syn1Neg.RowView(target)
						x := mat64.Dot(l1, l2)
						networkOutput := float64(1.0) / (1.0 + math.Exp(-x))

						g := (float64(label) - networkOutput) * model.Alpha

						neu1eSub := mat64.NewVector(model.VecDim, nil)
						neu1eSub.ScaleVec(g, l2)
						neu1e.AddVec(neu1e, neu1eSub)

						//if word.Literal == "computer" {
						//fmt.Printf("context word: %v\n", lastWord.Literal)
						//fmt.Printf("negative: %v\n", model.Vocab[target].Literal)
						//fmt.Printf("x: %v\n", x)
						//fmt.Printf("x: %v\n", round.ToEven(x, 6))
						//
						//fmt.Printf("networkOutput: %v\n", networkOutput)
						//fmt.Printf("gradient: %v\n", g)
						//}

						errorHiddenLayerWeights := mat64.NewVector(model.VecDim, nil)
						errorHiddenLayerWeights.ScaleVec(g, l1)

						newOutputLayerWeights := mat64.NewVector(model.VecDim, nil)
						newOutputLayerWeights.AddVec(model.Syn1Neg.RowView(target), errorHiddenLayerWeights)
						model.Syn1Neg.SetRow(target, newOutputLayerWeights.RawVector().Data)
					}

					oldWeights := model.Syn0.RowView(lastWord.Id)
					newWeights := mat64.NewVector(model.VecDim, nil)
					newWeights.AddVec(oldWeights, neu1e)
					model.Syn0.SetRow(lastWord.Id, newWeights.RawVector().Data)
					lastWord.Updated++
				}
				sentencePosition += 1
			}
		}
	}

	fmt.Printf("skippedWords size: %d\n", len(skippedWords))
	fmt.Printf("wordCount: %d\n", wordCount)
	fmt.Printf("Vocab: %d\n", len(model.Vocab))
	fmt.Printf("RawVocab: %d\n", len(model.RawVocab))
	// fmt.Printf("skippedWords: %v\n", skippedWords)
	return model

}
