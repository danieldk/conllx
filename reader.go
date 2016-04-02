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
	tokens  Sentence
}

// NewReader creates a new CoNLL-X reader from a buffered I/O reader.
// The caller is responsible for closing the provided reader.
func NewReader(r *bufio.Reader) *Reader {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	return &Reader{
		scanner: scanner,
		eof:     false,
	}
}

func parseColumns(line string) ([10]string, int) {
	var columns [10]string

	for i := 0; i < 10; i++ {
		end := strings.IndexByte(line, byte('\t'))

		if end == -1 {
			if len(line) == 0 {
				return columns, i
			} else {
				columns[i] = line
				return columns, i + 1
			}
		}

		columns[i] = line[:end]
		line = line[end+1:]
	}

	return columns, 10
}

// ReadSentence returns the next sentence. If there is no more data
// that can be read, io.EOF is returned as the error.
//
// The returned Sentence slice is only valid until the next call of
// ReadSentence. If you need to retain a sentence accross calls,
// it is safe to make a copy.
func (r *Reader) ReadSentence() (sentence Sentence, err error) {
	r.tokens = r.tokens[:0]

	if r.eof {
		return nil, io.EOF
	}

	for r.scanner.Scan() {
		line := r.scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			if len(r.tokens) == 0 {
				continue
			}

			break
		}

		parts, partsLen := parseColumns(line)

		token, err := processToken(parts[:partsLen])
		if err != nil {
			return nil, err
		}

		r.tokens = append(r.tokens, token)
	}

	if r.scanner.Err() == io.EOF {
		r.eof = true
	}

	if len(r.tokens) == 0 {
		return nil, io.EOF
	}

	return r.tokens, nil
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
