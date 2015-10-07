package conllx

import (
	"testing"
)

var projectiveSent1 []Token

var projectiveSent2 []Token

func init() {
	projectiveSent1 = make([]Token, 10)
	projectiveSent1[0].SetForm("Für").SetHead(4).SetHeadRel("PP|OBJA")
	projectiveSent1[1].SetForm("diese").SetHead(3).SetHeadRel("DET")
	projectiveSent1[2].SetForm("Behauptung").SetHead(1).SetHeadRel("PN")
	projectiveSent1[3].SetForm("hat").SetHead(0).SetHeadRel("ROOT")
	projectiveSent1[4].SetForm("Beckmeyer").SetHead(4).SetHeadRel("SUBJ")
	projectiveSent1[5].SetForm("bisher").SetHead(9).SetHeadRel("ADV")
	projectiveSent1[6].SetForm("keinen").SetHead(8).SetHeadRel("DET")
	projectiveSent1[7].SetForm("Nachweis").SetHead(9).SetHeadRel("OBJA")
	projectiveSent1[8].SetForm("geliefert").SetHead(4).SetHeadRel("AUX")
	projectiveSent1[9].SetForm(".").SetHead(9).SetHeadRel("-PUNCT-")

	projectiveSent2 = make([]Token, 11)
	projectiveSent2[0].SetForm("Auch").SetHead(2).SetHeadRel("ADV")
	projectiveSent2[1].SetForm("für").SetHead(5).SetHeadRel("PP|PN")
	projectiveSent2[2].SetForm("Rumänien").SetHead(2).SetHeadRel("PN")
	projectiveSent2[3].SetForm("selbst").SetHead(3).SetHeadRel("ADV")
	projectiveSent2[4].SetForm("ist").SetHead(0).SetHeadRel("ROOT")
	projectiveSent2[5].SetForm("der").SetHead(7).SetHeadRel("DET")
	projectiveSent2[6].SetForm("Papst-Besuch").SetHead(5).SetHeadRel("SUBJ")
	projectiveSent2[7].SetForm("von").SetHead(5).SetHeadRel("PRED")
	projectiveSent2[8].SetForm("großer").SetHead(10).SetHeadRel("ATTR")
	projectiveSent2[9].SetForm("Bedeutung").SetHead(8).SetHeadRel("PN")
	projectiveSent2[10].SetForm(".").SetHead(10).SetHeadRel("-PUNCT-")
}

func TestProjectivizeHead(t *testing.T) {
	projectivizer := HeadProjectivizer{}
	projective := projectivizer.Projectivize(nonProjectiveSent1)
	testEquals(t, projectiveSent1, projective)

	// Shouldn't give a change
	projective = projectivizer.Projectivize(projective)
	testEquals(t, projectiveSent1, projective)

	projective = projectivizer.Projectivize(nonProjectiveSent2)
	testEquals(t, projectiveSent2, projective)
}

func TestDeprojectivizeHead(t *testing.T) {
	projectivizer := HeadProjectivizer{}
	nonProjective := projectivizer.Deprojectivize(projectiveSent1)
	testEquals(t, nonProjectiveSent1, nonProjective)

	nonProjective = projectivizer.Deprojectivize(projectiveSent2)
	testEquals(t, nonProjectiveSent2, nonProjective)
}
