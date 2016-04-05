// Copyright 2015 The conllx Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conllx

import "errors"

// A FoldSet contains fold numbers. This type is used with a SplittingReader
// to indicate from which folds sentences should be returned.
type FoldSet map[int]interface{}

var _ SentenceReader = &SplittingReader{}

// SplittingReader is a wrapper around a (CoNLL-X) Reader that splits the
// corpus into folds.
type SplittingReader struct {
	reader *Reader
	nFolds int
	folds  FoldSet
	count  int
}

// NewSplittingReader creates a SplittingReader, that splits the data in
// 'nFolds' folds. The reader returns the sentences that are in 'folds'.
func NewSplittingReader(reader *Reader, nFolds int, folds FoldSet) (*SplittingReader, error) {
	if nFolds < 1 {
		return nil, errors.New("The data should be 'splitted' in at least 1 fold.")
	}

	return &SplittingReader{
		reader: reader,
		nFolds: nFolds,
		folds:  folds,
		count:  -1,
	}, nil
}

// ReadSentence returns the next sentence that is in one of the folds
// requested from the SplittingReader.
func (r *SplittingReader) ReadSentence() (sentence Sentence, err error) {
	for {
		sentence, err := r.reader.ReadSentence()
		if err != nil {
			return sentence, err
		}

		r.count++

		if r.count == r.nFolds {
			r.count = 0
		}

		if _, ok := r.folds[r.count]; ok {
			return sentence, nil
		}
	}
}
