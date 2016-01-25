package conllx

import "testing"

var stringerTestToken = NewToken().
	SetForm("Test").
	SetLemma("test").
	SetCoarsePosTag("N").
	SetPosTag("NN").
	SetFeatures(map[string]string{"pos": "N"}).
	SetHead(0).
	SetHeadRel("ROOT").
	SetPHead(2).
	SetPHeadRel("PROOT")

var stringerTestCheck = "Test	test	N	NN	pos:N	0	ROOT	2	PROOT"
var stringerEmptyCheck = "_	_	_	_	_	_	_	_	_"

func TestToken(t *testing.T) {
	if stringerTestToken.String() != stringerTestCheck {
		t.Fatalf("Stringer error. Expected:\n%s\nGot\n%s", stringerTestCheck, stringerTestToken.String())
	}

	emptyToken := NewToken()
	if emptyToken.String() != stringerEmptyCheck {
		t.Fatalf("Stringer error. Expected:\n%s\nGot\n%s", stringerEmptyCheck, emptyToken.String())
	}
}
