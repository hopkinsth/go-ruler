package ruler

/*
This struct is the main format for rules or conditions in ruler-compatable libraries.
Here's a sample in JSON format:
	{
		"comparator": "eq",
		"path": "person.name",
		"value": "James"
	}

Valid comparators are: eq, neq, lt, lte, gt, gte, contains (regex), ncontains (!regex)

This struct is exported here so that you can include it in your own JSON encoding/decoding,
but go-ruler has a facility to help decode your rules from JSON into its own structs.
*/
type Rule struct {
	Comparator string
	Path       string
	Value      interface{}
}

/*
A RulerRule combines a single rule and a whole set of rules and is used
when building rules programmatically through Ruler's Rule() function.
It's not meant to be created directly.
*/
type RulerRule struct {
	*Ruler
	*Rule
}

// adds an equals condition
func (rf *RulerRule) Eq(value interface{}) *RulerRule {
	return rf.compare(eq, value)
}

// adds a not equals condition
func (rf *RulerRule) Neq(value interface{}) *RulerRule {
	return rf.compare(neq, value)
}

// adds a less than condition
func (rf *RulerRule) Lt(value interface{}) *RulerRule {
	return rf.compare(lt, value)
}

// adds a less than or equal condition
func (rf *RulerRule) Lte(value interface{}) *RulerRule {
	return rf.compare(lte, value)
}

// adds a greater than condition
func (rf *RulerRule) Gt(value interface{}) *RulerRule {
	return rf.compare(gt, value)
}

// adds a greater than or equal to condition
func (rf *RulerRule) Gte(value interface{}) *RulerRule {
	return rf.compare(gte, value)
}

// adds a matches (regex) condition
func (rf *RulerRule) Matches(value interface{}) *RulerRule {
	return rf.compare(matches, value)
}

// adds a not matches condition (ncontains, in the way this thinks of it)
func (rf *RulerRule) NotMatches(value interface{}) *RulerRule {
	return rf.compare(ncontains, value)
}

// Stops chaining for the current rule, allowing you
// to add rules for other properties
func (rf *RulerRule) End() *Ruler {
	return rf.Ruler
}

// comparator will either create a new ruler filter and add its filter
func (rf *RulerRule) compare(comp int, value interface{}) *RulerRule {
	var comparator string
	switch comp {
	case eq:
		comparator = "eq"
	case neq:
		comparator = "neq"
	case lt:
		comparator = "lt"
	case lte:
		comparator = "lte"
	case gt:
		comparator = "gt"
	case gte:
		comparator = "gte"
	case contains:
		comparator = "contains"
	case matches:
		comparator = "matches"
	case ncontains:
		comparator = "ncontains"
	}

	// if this thing has a comparator already, we need to make a new ruler filter
	if rf.Comparator != "" {
		rf = &RulerRule{
			rf.Ruler,
			&Rule{
				comparator,
				rf.Path,
				value,
			},
		}
		// attach the new filter to the ruler
		rf.Ruler.rules = append(rf.Ruler.rules, rf.Rule)
	} else {
		//if there is no comparator, we can just set things on the current filter
		rf.Comparator = comparator
		rf.Value = value
	}

	return rf
}
