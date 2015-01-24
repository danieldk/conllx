// The conllx package provides a reader (and soon writer)
// for the CoNNL-X format:
//
// http://ilk.uvt.nl/conll/
package conllx

const formBit = uint32(1 << 1)
const lemmaBit = uint32(1 << 2)
const coarsePosTagBit = uint32(1 << 3)
const posTagBit = uint32(1 << 4)
const featuresBit = uint32(1 << 5)
const headBit = uint32(1 << 6)
const headRelBit = uint32(1 << 7)
const pHeadBit = uint32(1 << 8)
const pHeadRelBit = uint32(1 << 9)

// A CONLL-X token
type Token struct {
	available    uint32
	form         string
	lemma        string
	coarsePosTag string
	posTag       string
	features     string
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

// Get the coarse-grained POS tag, the second tuple element is
// false if there is no tag stored in this token.
func (t *Token) PosTag() (string, bool) {
	return t.posTag, t.available&posTagBit != 0
}

// Get the features field, the second tuple element is
// false if there are no features stored in this token.
func (t *Token) Features() (string, bool) {
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
