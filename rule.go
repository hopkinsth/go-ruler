package ruler

import (
	"strings"
)

type Filter struct {
	Comparator string
	Path       string
	Value      interface{}
}

type Rule struct {
	filters []*Filter
}

func (r *Rule) Test(o map[string]interface{}) bool {
	for _, f := range r.filters {
		val := pluck(o, f.Path)

		switch f.Comparator {
		case "eq":
			return f.Value == val
		default:
			//should probably return an error or something
			//but this is good for now
			//if comparator is not implemented, return false
			return false
		}
	}

	return false
}

// given a map, pull a property from it at some deeply nested depth
// this reimplements JS `pluck` in go: https://github.com/gjohnson/pluck
func pluck(o map[string]interface{}, path string) interface{} {
	// support dots for now ebcause thats all we need
	parts := strings.Split(path, ".")

	if len(parts) == 1 && o[parts[0]] != nil {
		// if there is only one part, just return that property value
		return o[parts[0]]
	} else if len(parts) > 1 && o[parts[0]] != nil {
		var prev map[string]interface{}
		var ok bool
		if prev, ok = o[parts[0]].(map[string]interface{}); !ok {
			// not an object type! ...or a map, yeah, that.
			return nil
		}

		for i := 1; i < len(parts)-1; i += 1 {
			// we need to check the existence of another
			// map[string]interface for every property along the way
			cp := parts[i]

			if prev[cp] == nil {
				// didn't find the property, it's missing
				return nil
			}
			var ok bool
			if prev, ok = prev[cp].(map[string]interface{}); !ok {
				return nil
			}
		}

		if prev[parts[len(parts)-1]] != nil {
			return prev[parts[len(parts)-1]]
		} else {
			return nil
		}
	}

	return nil
}
