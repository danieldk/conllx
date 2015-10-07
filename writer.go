package conllx

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Writer struct {
	first         bool
	writer        io.Writer
	projectivizer Projectivizer
	direction     ProjectivizeDirection
}

// Create a new writer.
func NewWriter(w io.Writer) Writer {
	return Writer{
		first:         true,
		writer:        w,
		projectivizer: nil,
	}
}

// Set a projectivizer to deprojectivize dependency structures.
func (w *Writer) SetProjectivizer(projectivizer Projectivizer, direction ProjectivizeDirection) {
	w.projectivizer = projectivizer
	w.direction = direction
}

func (w *Writer) WriteSentence(sentence []Token) error {
	// Sentences are split using an extra newline. Moreover, there shouldn't
	// be a newline after the last token of the stream. So, we always print
	// the last token of the sentence without a newline and print two newlines
	// before each sentence (except for the first sentence).

	if w.first {
		w.first = false
	} else {
		fmt.Fprint(w.writer, "\n\n")
	}

	sentenceLen := len(sentence)

	if w.projectivizer != nil {
		switch w.direction {
		case Projectivize:
			sentence = w.projectivizer.Projectivize(sentence)
		case Deprojectivize:
			sentence = w.projectivizer.Deprojectivize(sentence)
		default:
			panic("Unknown projectivization direction.")
		}
	}

	for idx, token := range sentence {
		if idx == sentenceLen-1 {
			fmt.Fprintf(w.writer, "%d\t%s", idx+1, w.formatToken(token))
		} else {
			fmt.Fprintf(w.writer, "%d\t%s\n", idx+1, w.formatToken(token))
		}
	}

	return nil
}

func (w Writer) formatToken(token Token) string {
	cols := []string{
		w.formatColumn(token.Form),
		w.formatColumn(token.Lemma),
		w.formatColumn(token.CoarsePosTag),
		w.formatColumn(token.PosTag),
		w.formatFeaturesColumn(token.Features),
		w.formatIntColumn(token.Head),
		w.formatColumn(token.HeadRel),
		w.formatIntColumn(token.PHead),
		w.formatColumn(token.PHeadRel),
	}

	return strings.Join(cols, "\t")
}

type conllxGet func() (string, bool)
type conllxGetInt func() (uint, bool)
type conllxGetFeatures func() (*Features, bool)

func (w Writer) formatColumn(getter conllxGet) string {
	if value, ok := getter(); ok {
		return value
	} else {
		return "_"
	}
}

func (w Writer) formatIntColumn(getter conllxGetInt) string {
	if value, ok := getter(); ok {
		return strconv.FormatUint(uint64(value), 10)
	} else {
		return "_"
	}
}

func (w Writer) formatFeaturesColumn(getter conllxGetFeatures) string {
	if value, ok := getter(); ok {
		return value.FeaturesString()
	} else {
		return "_"
	}
}
