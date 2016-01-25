// Copyright 2015 The conllx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conllx

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// A Reader for CONLL-X files.
type Reader struct {
	scanner *bufio.Scanner
	eof     bool
}

// NewReader creates a new CoNLL-X reader from a buffered I/O reader.
// The caller is responsible for closing the provided reader.
func NewReader(r *bufio.Reader) Reader {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	return Reader{
		scanner: scanner,
		eof:     false,
	}
}

// ReadSentence returns the next sentence. If there is no more data
// that can be read, io.EOF is returned as the error.
func (r *Reader) ReadSentence() (sentence []Token, err error) {
	var tokens []Token

	if r.eof {
		return nil, io.EOF
	}

	for r.scanner.Scan() {
		line := r.scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			if len(tokens) == 0 {
				continue
			}

			break
		}

		parts := strings.Split(line, "\t")
		token, err := processToken(parts)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	if r.scanner.Err() == io.EOF {
		r.eof = true
	}

	if len(tokens) == 0 {
		return nil, io.EOF
	}

	return tokens, nil
}

func processToken(columns []string) (Token, error) {
	_, _, err := intValueForColumn(columns, 0)
	if err != nil {
		return Token{}, err
	}

	form, formBit := valueForColumn(columns, 1)
	lemma, lemmaBit := valueForColumn(columns, 2)
	cTag, cTagBit := valueForColumn(columns, 3)
	tag, tagBit := valueForColumn(columns, 4)
	features, featuresBit := valueForColumn(columns, 5)
	headRel, headRelBit := valueForColumn(columns, 7)
	pHeadRel, pHeadRelBit := valueForColumn(columns, 9)

	head, headBit, err := intValueForColumn(columns, 6)
	if err != nil {
		return Token{}, err
	}

	pHead, pHeadBit, err := intValueForColumn(columns, 8)
	if err != nil {
		return Token{}, err
	}

	var featuresField *Features
	if featuresBit != 0 {
		featuresField = newFeatures(features)
	}

	return Token{
		available: formBit | lemmaBit | cTagBit | tagBit | featuresBit |
			headBit | headRelBit | pHeadBit | pHeadRelBit,
		form:         form,
		lemma:        lemma,
		coarsePosTag: cTag,
		posTag:       tag,
		features:     featuresField,
		head:         head,
		headRel:      headRel,
		pHead:        pHead,
		pHeadRel:     pHeadRel,
	}, nil
}

// Return the value for a column, returns the corresponding bit
// set to one if the value was actually present.
func valueForColumn(columns []string, idx int) (string, fields) {
	if idx >= len(columns) || columns[idx] == "_" {
		return "", 0
	}

	return columns[idx], fields(1) << fields(idx-1)
}

// Return the value for a column, returns the corresponding bit
// set to one if the value was actually present.
func intValueForColumn(columns []string, idx int) (uint, fields, error) {
	if idx >= len(columns) || columns[idx] == "_" {
		return 0, 0, nil
	}

	val, err := strconv.ParseUint(columns[idx], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return uint(val), fields(1) << fields(idx-1), nil
}
