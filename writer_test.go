// Copyright 2015 The conllx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conllx

import (
	"bytes"
	"fmt"
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

func ExampleWriter() {
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	writer.WriteSentence(
		Sentence{
			*NewToken().SetForm("Hello").SetPosTag("expr"),
			*NewToken().SetForm("world").SetPosTag("noun")})
	writer.WriteSentence(
		Sentence{
			*NewToken().SetForm("Go").SetPosTag("name"),
			*NewToken().SetForm("rocks").SetPosTag("verb")})

	fmt.Println(buf.String())

	// Output:
	// 1	Hello	_	_	expr	_	_	_	_	_
	// 2	world	_	_	noun	_	_	_	_	_
	//
	// 1	Go	_	_	name	_	_	_	_	_
	// 2	rocks	_	_	verb	_	_	_	_	_
}
