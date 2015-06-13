package ruler

type Rule struct {
	Comparator string
	Path       string
	Value      interface{}
}

// this is a special struct used for
// implementing the programmatic rule building
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

// adds a greater than condition
func (rf *RulerRule) NotMatches(value interface{}) *RulerRule {
	return rf.compare(ncontains, value)
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
