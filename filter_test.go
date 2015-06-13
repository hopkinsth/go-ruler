package ruler

import "testing"

func TestNewFilterWhenComparatorExists(t *testing.T) {
	rule := &Rule{
		Comparator: "eq",
		Path:       "name",
		Value:      "Bob",
	}

	r := &Ruler{
		rules: []*Rule{rule},
	}

	rf := &RulerRule{
		r,
		rule,
	}

	res := rf.Gt("sup")

	if res == rf {
		t.Error("RulerFilter should return a new struct if there is an existing comparator!")
	}
}
