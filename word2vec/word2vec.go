package word2vec

import (
	"fmt"
	"github.com/gonum/matrix"
)

// Trains a word2vec model after Mikolov: https://arxiv.org/abs/1310.4546

// python implementation: https://github.com/RaRe-Technologies/gensim/blob/develop/gensim/models/word2vec.py

// Model contains the word2vec vocab, vectors, parameters, etc.
type Model struct {
	Vocab map[string]Phrase
}

// Model contains the word2vec vocab, vectors, parameters, etc.
type Phrase struct {
	Literal string
	Vector  float64
	Count   float64
}

func train_sg_pair()
