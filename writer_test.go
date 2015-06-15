package conllx

import (
	"bytes"
	"testing"
)

const testFragmentExplicit string = `1	Die	die	ART	ART	nsf	2	DET	_	_
2	Großaufnahme	Großaufnahme	N	NN	nsf	0	ROOT	_	_

1	Gilles	Gilles	N	NE	nsm	0	ROOT	_	_
2	Deleuze	Deleuze	N	NE	case:nominative|number:singular|gender:masculine	1	APP	_	_`

func TestWriter(t *testing.T) {
	sentences := [][]Token{
		testFragmentSent1,
		testFragmentSent2,
	}

	writerTestHelper(t, sentences)
}

func writerTestHelper(t *testing.T, sentences [][]Token) {
	var buffer bytes.Buffer
	writer := NewWriter(&buffer)

	for _, sentence := range sentences {
		writer.WriteSentence(sentence)
	}

	if buffer.String() != testFragmentExplicit {
		t.Fatalf("Got:\n%sExpected:\n%s\n", buffer.String(), testFragmentExplicit)
	}
}
