package graphb

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	r := MakeRequest().SetElements(
		MakeFragment("alias", "table").
			SetFields(
				MakeField("testField"),
			),
		MakeQuery(TypeQuery).
			SetFields(
				MakeFragmentField("alias"),
			),
	)
	fmt.Println(r.JSON())
}
