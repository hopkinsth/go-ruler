package ruler

import "testing"

func TestNewFilterWhenComparatorExists(t *testing.T) {
	filter := &Filter{
		Comparator: "eq",
		Path:       "name",
		Value:      "Bob",
	}

	r := &Ruler{
		filters: []*Filter{filter},
	}

	rf := &RulerFilter{
		r,
		filter,
	}

	res := rf.Gt("sup")

	if res == rf {
		t.Error("RulerFilter should return a new struct if there is an existing comparator!")
	}
}
