package conllx

import (
	"testing"
)

var nonProjectiveSent1 []Token
var nonProjectiveSent2 []Token

var nonProjectiveEdges1 = []edge{
	edge{8, 1, "PP"},
}
var nonProjectiveEdges2 = []edge{
	edge{10, 2, "PP"},
}

func init() {
	nonProjectiveSent1 = make([]Token, 10)
	nonProjectiveSent1[0].SetForm("Für").SetHead(8).SetHeadRel("PP")
	nonProjectiveSent1[1].SetForm("diese").SetHead(3).SetHeadRel("DET")
	nonProjectiveSent1[2].SetForm("Behauptung").SetHead(1).SetHeadRel("PN")
	nonProjectiveSent1[3].SetForm("hat").SetHead(0).SetHeadRel("ROOT")
	nonProjectiveSent1[4].SetForm("Beckmeyer").SetHead(4).SetHeadRel("SUBJ")
	nonProjectiveSent1[5].SetForm("bisher").SetHead(9).SetHeadRel("ADV")
	nonProjectiveSent1[6].SetForm("keinen").SetHead(8).SetHeadRel("DET")
	nonProjectiveSent1[7].SetForm("Nachweis").SetHead(9).SetHeadRel("OBJA")
	nonProjectiveSent1[8].SetForm("geliefert").SetHead(4).SetHeadRel("AUX")
	nonProjectiveSent1[9].SetForm(".").SetHead(9).SetHeadRel("-PUNCT-")

	nonProjectiveSent2 = make([]Token, 11)
	nonProjectiveSent2[0].SetForm("Auch").SetHead(2).SetHeadRel("ADV")
	nonProjectiveSent2[1].SetForm("für").SetHead(10).SetHeadRel("PP")
	nonProjectiveSent2[2].SetForm("Rumänien").SetHead(2).SetHeadRel("PN")
	nonProjectiveSent2[3].SetForm("selbst").SetHead(3).SetHeadRel("ADV")
	nonProjectiveSent2[4].SetForm("ist").SetHead(0).SetHeadRel("ROOT")
	nonProjectiveSent2[5].SetForm("der").SetHead(7).SetHeadRel("DET")
	nonProjectiveSent2[6].SetForm("Papst-Besuch").SetHead(5).SetHeadRel("SUBJ")
	nonProjectiveSent2[7].SetForm("von").SetHead(5).SetHeadRel("PRED")
	nonProjectiveSent2[8].SetForm("großer").SetHead(10).SetHeadRel("ATTR")
	nonProjectiveSent2[9].SetForm("Bedeutung").SetHead(8).SetHeadRel("PN")
	nonProjectiveSent2[10].SetForm(".").SetHead(10).SetHeadRel("-PUNCT-")
}

func TestNonProjectiveArcs(t *testing.T) {
	g := sentenceToDepGraph(nonProjectiveSent1)
	np := nonProjectiveArcs(g)
	testEquals(t, nonProjectiveEdges1, np)

	g = sentenceToDepGraph(nonProjectiveSent2)
	np = nonProjectiveArcs(g)
	testEquals(t, nonProjectiveEdges2, np)
}
