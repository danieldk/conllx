// The conllx package provides a reader (and soon writer)
// for the CoNNL-X format:
//
// http://ilk.uvt.nl/conll/
package conllx

import (
	"fmt"
	"strings"
)

type fields uint32

const (
	formBit fields = 1 << iota
	lemmaBit
	coarsePosTagBit
	posTagBit
	featuresBit
	headBit
	headRelBit
	pHeadBit
	pHeadRelBit
)

// Features from the CONLL-X features field.
type Features struct {
	featuresString string
	featuresMap    map[string]string
}

// Construct a new features field from a features string.
func newFeatures(featuresString string) *Features {
	return &Features{
		featuresString: featuresString,
		featuresMap:    nil,
	}
}

// Get the features as a string. This will give feature in
// exactly the same format as the original CONLL-X data.
func (f *Features) FeaturesString() string {
	return f.featuresString
}

// Get the features as a key-value mapping. Features that do
// not follow the expected format are skipped.
//
// The feature map is lazily initialized on its first call. No
// feature field parsing is done if this method is not called.
func (f *Features) FeaturesMap() map[string]string {
	if f.featuresMap == nil {
		f.featuresMap = make(map[string]string)
		for _, av := range strings.Split(f.featuresString, "|") {
			if sepIdx := strings.IndexByte(av, ':'); sepIdx != -1 {
				f.featuresMap[av[:sepIdx]] = av[sepIdx+1:]
			}
		}
	}

	return f.featuresMap
}

// A CONLL-X token
type Token struct {
	available    fields
	form         string
	lemma        string
	coarsePosTag string
	posTag       string
	features     *Features
	head         uint
	headRel      string
	pHead        uint
	pHeadRel     string
}

// Get the form, the second tuple element is false if there
// is no form stored in this token.
func (t *Token) Form() (string, bool) {
	return t.form, t.available&formBit != 0
}

// Get the lemma, the second tuple element is false if there
// is no lemma stored in this token.
func (t *Token) Lemma() (string, bool) {
	return t.lemma, t.available&lemmaBit != 0
}

// Get the coarse-grained POS tag, the second tuple element is
// false if there is no coarse-grained tag stored in this token.
func (t *Token) CoarsePosTag() (string, bool) {
	return t.coarsePosTag, t.available&coarsePosTagBit != 0
}

// Get the fine-grained POS tag, the second tuple element is
// false if there is no tag stored in this token.
func (t *Token) PosTag() (string, bool) {
	return t.posTag, t.available&posTagBit != 0
}

// Get the features field, the second tuple element is
// false if there are no features stored in this token.
func (t *Token) Features() (*Features, bool) {
	return t.features, t.available&featuresBit != 0
}

// Get the head field, the second tuple element is
// false if there is no head stored in this token.
func (t *Token) Head() (uint, bool) {
	return t.head, t.available&headBit != 0
}

// Get the head relation, the second tuple element is
// false if there is no head relation stored in this token.
func (t *Token) HeadRel() (string, bool) {
	return t.headRel, t.available&headRelBit != 0
}

// Get the projective head field, the second tuple element is
// false if there is no projective head stored in this token.
func (t *Token) PHead() (uint, bool) {
	return t.pHead, t.available&pHeadBit != 0
}

// Get the head relation, the second tuple element is
// false if there is no head relation stored in this token.
func (t *Token) PHeadRel() (string, bool) {
	return t.pHeadRel, t.available&pHeadRelBit != 0
}

// Set the features for this token.
func (t *Token) SetFeatures(features map[string]string) {
	f := new(Features)
	f.featuresMap = features

	fVals := make([]string, 0, len(features))
	for k, v := range features {
		fVals = append(fVals, fmt.Sprintf("%s:%s", k, v))
	}

	f.featuresString = strings.Join(fVals, "|")

	t.features = f
	t.available |= featuresBit
}

// Set the form field.
func (t *Token) SetForm(form string) *Token {
	t.form = form
	t.available |= formBit
	return t
}

// Set the lemma field.
func (t *Token) SetLemma(lemma string) *Token {
	t.lemma = lemma
	t.available |= lemmaBit
	return t
}

// Set the coarse-grained POS tag field.
func (t *Token) SetCoarsePosTag(coarsePosTag string) *Token {
	t.coarsePosTag = coarsePosTag
	t.available |= coarsePosTagBit
	return t
}

// Set the POS tag field.
func (t *Token) SetPosTag(posTag string) *Token {
	t.posTag = posTag
	t.available |= posTagBit
	return t
}

// Set the head field.
func (t *Token) SetHead(head uint) *Token {
	t.head = head
	t.available |= headBit
	return t
}

// Set the head relation.
func (t *Token) SetHeadRel(rel string) *Token {
	t.headRel = rel
	t.available |= headRelBit
	return t
}

// Set the projective head field.
func (t *Token) SetPHead(head uint) *Token {
	t.pHead = head
	t.available |= pHeadBit
	return t
}

// Set the projective head relation field.
func (t *Token) SetPHeadRel(rel string) *Token {
	t.pHeadRel = rel
	t.available |= pHeadRelBit
	return t
}
