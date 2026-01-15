package graphb

import (
	"fmt"
	"testing"
)

func TestFragment(t *testing.T) {
	f := MakeFragment("alias", "table").
		SetFields(
			MakeField("testField"),
		)
	fmt.Println(f.JSON())
}
