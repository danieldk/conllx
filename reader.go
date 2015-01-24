package conll

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type CONLLReader struct {
	reader *bufio.Reader
	eof    bool
}

func NewCONLLReader(r *bufio.Reader) CONLLReader {
	return CONLLReader{r, false}
}

func (r *CONLLReader) ReadSentence() (sentence []Token, err error) {
	tokens := make([]Token, 0)

	if r.eof {
		return nil, io.EOF
	}

	for !r.eof {
		line, err := r.reader.ReadString('\n')

		// If we have an error and it is EOF, we still want to process
		// the last line.
		if err != nil {
			if err == io.EOF {
				r.eof = true
			} else {
				return nil, err
			}
		}

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

	return tokens, nil
}

func processToken(columns []string) (Token, error) {
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

	return Token{
		available: formBit | lemmaBit | cTagBit | tagBit | featuresBit |
			headBit | headRelBit | pHeadBit | pHeadRelBit,
		form:         form,
		lemma:        lemma,
		coarsePosTag: cTag,
		posTag:       tag,
		features:     features,
		head:         head,
		headRel:      headRel,
		pHead:        pHead,
		pHeadRel:     pHeadRel,
	}, nil
}

// Return the value for a column, returns the corresponding bit
// set to one if the value was actually present.
func valueForColumn(columns []string, idx int) (string, uint32) {
	if idx >= len(columns) || columns[idx] == "_" {
		return "", 0
	}

	return columns[idx], uint32(1) << uint32(idx)
}

// Return the value for a column, returns the corresponding bit
// set to one if the value was actually present.
func intValueForColumn(columns []string, idx int) (uint, uint32, error) {
	if idx >= len(columns) || columns[idx] == "_" {
		return 0, 0, nil
	}

	val, err := strconv.ParseUint(columns[idx], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	return uint(val), uint32(1) << uint32(idx), nil
}
