package conllx

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

const testFragment string = `1	Die	die	ART	ART	nsf	2	DET
2	Großaufnahme	Großaufnahme	N	NN	nsf	0	ROOT

1	Gilles	Gilles	N	NE	nsm	0	ROOT
2	Deleuze	Deleuze	N	NE	nsm	1	APP`

// Not according to CoNLL-X, but we want to handle it anyway.
const testFragmentRobust string = `1	Die	die	ART	ART	nsf	2	DET
2	Großaufnahme	Großaufnahme	N	NN	nsf	0	ROOT


1	Gilles	Gilles	N	NE	nsm	0	ROOT
2	Deleuze	Deleuze	N	NE	nsm	1	APP`

const testFragmentMarkedEmpty string = `1	Die	die	ART	ART	nsf	2	DET	_	_
2	Großaufnahme	Großaufnahme	N	NN	nsf	0	ROOT	_	_

1	Gilles	Gilles	N	NE	nsm	0	ROOT	_	_
2	Deleuze	Deleuze	N	NE	nsm	1	APP	_	_`

var testFragmentSent1 []Token = []Token{
	Token{0xFE, "Die", "die", "ART", "ART", "nsf", 2, "DET", 0, ""},
	Token{0xFE, "Großaufnahme", "Großaufnahme", "N", "NN", "nsf", 0, "ROOT", 0, ""},
}

var testFragmentSent2 []Token = []Token{
	Token{0xFE, "Gilles", "Gilles", "N", "NE", "nsm", 0, "ROOT", 0, ""},
	Token{0xFE, "Deleuze", "Deleuze", "N", "NE", "nsm", 1, "APP", 0, ""},
}

func equalOrFail(t *testing.T, correct, test []Token) {
	if !reflect.DeepEqual(correct, test) {
		t.Fatalf("Parsing error:\nCorrect: %v\nGot: %v", correct, test)
	}
}

func testHelper(t *testing.T, sentenceString string) {
	r := stringReader(sentenceString)

	sentence, err := r.ReadSentence()
	if err != nil {
		t.Fatalf("Sentence read should succeed: %s", err)
	}

	equalOrFail(t, testFragmentSent1, sentence)

	sentence2, err := r.ReadSentence()
	if err != nil {
		t.Fatalf("Sentence read should succeed: %s", err)
	}

	equalOrFail(t, testFragmentSent2, sentence2)

	_, err = r.ReadSentence()
	if err != io.EOF {
		t.Fatalf("Reader should return EOF.")
	}
}

func TestCorrect(t *testing.T) {
	testHelper(t, testFragment)
}

func TestCorrectRobust(t *testing.T) {
	testHelper(t, testFragmentRobust)
}

func TestCorrectMarkedEmpty(t *testing.T) {
	testHelper(t, testFragmentMarkedEmpty)
}

func TestEmpty(t *testing.T) {
	r := stringReader("")
	if _, err := r.ReadSentence(); err != io.EOF {
		t.Fatal("Parsing the empty string should return EOF")
	}

	r = stringReader("\n\n\n\n")
	if _, err := r.ReadSentence(); err != io.EOF {
		t.Fatal("Parsing the empty string should return EOF")
	}
}

func stringReader(s string) Reader {
	reader := strings.NewReader(s)
	return NewReader(bufio.NewReader(reader))
}
