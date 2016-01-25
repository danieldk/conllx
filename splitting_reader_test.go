package conllx

import (
	"io"
	"testing"
)

const splitTestFragment string = `1	a
2	b

1	c
2	d

1	e
2	f

1	g
2	h

1	i
2	j`

var splitTestSent0 = []Token{
	*NewToken().SetForm("a"),
	*NewToken().SetForm("b"),
}

var splitTestSent2 = []Token{
	*NewToken().SetForm("e"),
	*NewToken().SetForm("f"),
}

var splitTestSent3 = []Token{
	*NewToken().SetForm("g"),
	*NewToken().SetForm("h"),
}

var splitTestSent4 = []Token{
	*NewToken().SetForm("i"),
	*NewToken().SetForm("j"),
}

type splittingTestCase struct {
	nFolds   int
	folds    FoldSet
	expected [][]Token
}

var splittingTestCases = []splittingTestCase{
	{
		nFolds: 2,
		folds:  map[int]interface{}{0: nil},
		expected: [][]Token{
			splitTestSent0,
			splitTestSent2,
			splitTestSent4,
		},
	},
	{
		nFolds: 3,
		folds:  map[int]interface{}{0: nil, 2: nil},
		expected: [][]Token{
			splitTestSent0,
			splitTestSent2,
			splitTestSent3,
		},
	},
}

func TestSplitting(t *testing.T) {
	for _, testCase := range splittingTestCases {
		r, err := NewSplittingReader(stringReader(splitTestFragment), testCase.nFolds, testCase.folds)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		for _, correct := range testCase.expected {
			sent, err := r.ReadSentence()
			equalOrFail(t, err, correct, sent)
		}

		_, err = r.ReadSentence()
		if err != io.EOF {
			t.Fatalf("Reader should return EOF.")
		}
	}
}

func TestSplittingNoFolds(t *testing.T) {
	folds := map[int]interface{}{}
	r, err := NewSplittingReader(stringReader(splitTestFragment), 2, folds)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	_, err = r.ReadSentence()
	if err != io.EOF {
		t.Fatal("expected immediate EOF when reading no folds")
	}
}
