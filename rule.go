package ruler

import (
	_ "reflect"
	_ "regexp"
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

// func (r *Rule) Test(o map[string]interface{}) bool {
// 	for _, f := range r.filters {
// 		val := pluck(o, f.Path)
// 		rtype := reflect.TypeOf(f.Value)
// 		atype := reflect.Type(val)

// 		if rtype.Name() != atype.Name() {

// 		}
// 	}
// }

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

		prev = o[parts[0]].(map[string]interface{})

		for i := 1; i < len(parts)-1; i += 1 {
			// we need to check the existence of another
			// map[string]interface for every property along the way
			cp := parts[i]

			if prev[cp] == nil {
				// didn't find the property, it's missing
				return nil
			}

			prev = prev[cp].(map[string]interface{})
		}

		if prev[parts[len(parts)-1]] != nil {
			return prev[parts[len(parts)-1]]
		} else {
			return nil
		}
	}

	return nil
}
