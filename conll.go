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

// FeaturesString returns the token features as a string. This will give
// feature in exactly the same format as the original CONLL-X data.
func (f *Features) FeaturesString() string {
	return f.featuresString
}

// FeaturesMap returns the token features as a key-value mapping. Features
// that do not follow the expected format are skipped.
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

// Token stores a token with the CONLL-X annotation layers.
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

// Form returns the form (the actual token), the second tuple element is
// false when there is no form stored in this token.
func (t *Token) Form() (string, bool) {
	return t.form, t.available&formBit != 0
}

// Lemma returns the lemma of the token, the second tuple element is false
// when there is no lemma stored in this token.
func (t *Token) Lemma() (string, bool) {
	return t.lemma, t.available&lemmaBit != 0
}

// CoarsePosTag returns the coarse-grained POS tag of the token, the
// second tuple element is false when there is no coarse-grained tag
// stored in this token.
func (t *Token) CoarsePosTag() (string, bool) {
	return t.coarsePosTag, t.available&coarsePosTagBit != 0
}

// PosTag returns the fine-grained POS tag of the token, the second
// tuple element is false when there is no fine-grained tag stored in
// this token.
func (t *Token) PosTag() (string, bool) {
	return t.posTag, t.available&posTagBit != 0
}

// Features returns the features field, the second tuple element is false
// when there are no features stored in this token.
func (t *Token) Features() (*Features, bool) {
	return t.features, t.available&featuresBit != 0
}

// Head returns the head of the token, the second tuple element is false
// when there is no head stored in this token.
func (t *Token) Head() (uint, bool) {
	return t.head, t.available&headBit != 0
}

// HeadRel returns the relation of the token to its head, the second
// tuple element is false when there is no head relation stored in this
// token.
func (t *Token) HeadRel() (string, bool) {
	return t.headRel, t.available&headRelBit != 0
}

// PHead returns the projective head of the token, the second tuple
// element is false when there is no head stored in this token.
func (t *Token) PHead() (uint, bool) {
	return t.pHead, t.available&pHeadBit != 0
}

// PHeadRel returns the relation of the token to its projective head, the
// second tuple element is false when there is no head relation stored in
// this token.
func (t *Token) PHeadRel() (string, bool) {
	return t.pHeadRel, t.available&pHeadRelBit != 0
}

// SetFeatures sets the features for this token. The token itself is
// returned to allow method chaining.
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

// SetForm sets the form for this token. The token itself is returned to
// allow method chaining.
func (t *Token) SetForm(form string) *Token {
	t.form = form
	t.available |= formBit
	return t
}

// SetLemma sets the lemma for this token. The token itself is returned to
// allow method chaining.
func (t *Token) SetLemma(lemma string) *Token {
	t.lemma = lemma
	t.available |= lemmaBit
	return t
}

// SetCoarsePosTag sets the coarse-grained POS tag for this token. The
// token itself is returned to allow method chaining.
func (t *Token) SetCoarsePosTag(coarsePosTag string) *Token {
	t.coarsePosTag = coarsePosTag
	t.available |= coarsePosTagBit
	return t
}

// SetPosTag sets the fine-grained POS tag for this token. The token
// itself is returned to allow method chaining.
func (t *Token) SetPosTag(posTag string) *Token {
	t.posTag = posTag
	t.available |= posTagBit
	return t
}

// SetHead sets the head of this token. The token itself is returned to
// allow method chaining.
func (t *Token) SetHead(head uint) *Token {
	t.head = head
	t.available |= headBit
	return t
}

// SetHeadRel sets the relation to the head of this token. The token
// itself is returned to allow method chaining.
func (t *Token) SetHeadRel(rel string) *Token {
	t.headRel = rel
	t.available |= headRelBit
	return t
}

// SetPHead sets the projective head of this token. The token itself is
// returned to allow method chaining.
func (t *Token) SetPHead(head uint) *Token {
	t.pHead = head
	t.available |= pHeadBit
	return t
}

// SetPHeadRel sets the relation to the projective head of this token.
// The token itself is returned to allow method chaining.
func (t *Token) SetPHeadRel(rel string) *Token {
	t.pHeadRel = rel
	t.available |= pHeadRelBit
	return t
}
