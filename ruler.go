package ruler

import (
	"encoding/json"
	"github.com/tj/go-debug"
	"reflect"
	"regexp"
	"strings"
)

var ruleDebug = debug.Debug("ruler:rule")

// we'll use these values
// to avoid passing strings to our
// special comparison func for these comparators
const (
	eq        = iota
	neq       = iota
	gt        = iota
	gte       = iota
	lt        = iota
	lte       = iota
	exists    = iota
	nexists   = iota
	regex     = iota
	matches   = iota
	contains  = iota
	ncontains = iota
)

type Ruler struct {
	filters []*Filter
}

func NewRuler(filters *[]*Filter) *Ruler {
	if filters != nil {
		return &Ruler{
			*filters,
		}
	}

	return &Ruler{}
}

func NewRulerWithJson(jsonstr []byte) (*Ruler, error) {
	var filters *[]*Filter

	err := json.Unmarshal(jsonstr, &filters)
	if err != nil {
		return nil, err
	}

	return NewRuler(filters), nil
}

// begins chaining of rules
func (r *Ruler) Rule(path string) *RulerFilter {
	filter := &Filter{
		"",
		path,
		nil,
	}

	return &RulerFilter{
		r,
		filter,
	}
}

func (r *Ruler) Test(o map[string]interface{}) bool {
	for _, f := range r.filters {
		val := pluck(o, f.Path)

		if val != nil {
			// both the actual and expected value must be comparable
			a := reflect.TypeOf(val)
			e := reflect.TypeOf(f.Value)

			if !a.Comparable() || !e.Comparable() {
				return false
			}

			if !r.Compare(f, val) {
				return false
			}
		} else if val == nil && (f.Comparator == "exists" || f.Comparator == "nexists") {
			// either one of these can be done
			return r.Compare(f, val)
		} else {
			ruleDebug("did not find property (%s) on map", f.Path)
			// if we couldn't find the value on the map
			// and the comparator isn't exists/nexists, this fails
			return false
		}

	}

	return true
}

// compares real v. actual values
func (r *Ruler) Compare(f *Filter, actual interface{}) bool {
	ruleDebug("beginning comparison")
	expected := f.Value
	switch f.Comparator {
	case "eq":
		return actual == expected

	case "neq":
		return actual != expected

	case "gt":
		return r.CompareInequality(gt, actual, expected)

	case "gte":
		return r.CompareInequality(gte, actual, expected)

	case "lt":
		return r.CompareInequality(lt, actual, expected)

	case "lte":
		return r.CompareInequality(lte, actual, expected)

	case "exists":
		// not sure this makes complete sense
		return actual != nil

	case "nexists":
		return actual == nil

	case "regex":
		fallthrough
	case "contains":
		fallthrough
	case "matches":
		return r.CompareRegexp(actual, expected)

	case "ncontains":
		return !r.CompareRegexp(actual, expected)
	default:
		//should probably return an error or something
		//but this is good for now
		//if comparator is not implemented, return false
		return false
	}
}

// runs equality comparison
// separated in a different function because
// we need to do another type assertion here
// and some other acrobatics
func (r *Ruler) CompareInequality(op int, actual, expected interface{}) bool {
	// need some variables for these deals
	ruleDebug("entered inequality comparison")
	var cmpStr [2]string
	var cmpUint [2]uint64
	var cmpInt [2]int64
	var cmpFloat [2]float64

	for idx, i := range []interface{}{actual, expected} {
		switch t := i.(type) {
		case uint8:
			cmpUint[idx] = uint64(t)
		case uint16:
			cmpUint[idx] = uint64(t)
		case uint32:
			cmpUint[idx] = uint64(t)
		case uint64:
			cmpUint[idx] = t
		case uint:
			cmpUint[idx] = uint64(t)
		case int8:
			cmpInt[idx] = int64(t)
		case int16:
			cmpInt[idx] = int64(t)
		case int32:
			cmpInt[idx] = int64(t)
		case int64:
			cmpInt[idx] = t
		case int:
			cmpInt[idx] = int64(t)
		case float32:
			cmpFloat[idx] = float64(t)
		case float64:
			cmpFloat[idx] = t
		case string:
			cmpStr[idx] = t
		default:
			ruleDebug("invalid type for inequality comparison")
			return false
		}
	}

	// whichever of these works, we're happy with
	// but if you're trying to compare a string to an int, oh well!
	switch op {
	case gt:
		return cmpStr[0] > cmpStr[1] ||
			cmpUint[0] > cmpUint[1] ||
			cmpInt[0] > cmpInt[1] ||
			cmpFloat[0] > cmpFloat[1]
	case gte:
		return cmpStr[0] >= cmpStr[1] ||
			cmpUint[0] >= cmpUint[1] ||
			cmpInt[0] >= cmpInt[1] ||
			cmpFloat[0] >= cmpFloat[1]
	case lt:
		return cmpStr[0] < cmpStr[1] ||
			cmpUint[0] < cmpUint[1] ||
			cmpInt[0] < cmpInt[1] ||
			cmpFloat[0] < cmpFloat[1]
	case lte:
		return cmpStr[0] <= cmpStr[1] ||
			cmpUint[0] <= cmpUint[1] ||
			cmpInt[0] <= cmpInt[1] ||
			cmpFloat[0] <= cmpFloat[1]
	}

	return false
}

func (r *Ruler) CompareRegexp(actual, expected interface{}) bool {
	ruleDebug("beginning regexp")
	// regexps must be strings
	var streg string
	var ok bool
	if streg, ok = expected.(string); !ok {
		ruleDebug("expected value not actually a string, bailing")
		return false
	}

	var astring string
	if astring, ok = actual.(string); !ok {
		ruleDebug("actual value not actually a string, bailing")
		return false
	}

	reg, err := regexp.Compile(streg)
	if err != nil {
		ruleDebug("regexp is bad, bailing")
		return false
	}

	return reg.MatchString(astring)
}

// given a map, pull a property from it at some deeply nested depth
// this reimplements (most of) JS `pluck` in go: https://github.com/gjohnson/pluck
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
