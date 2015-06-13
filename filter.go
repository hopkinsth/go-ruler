package ruler

type Filter struct {
	Comparator string
	Path       string
	Value      interface{}
}

// this is a special struct used for
// implementing the programmatic rule building
type RulerFilter struct {
	*Ruler
	*Filter
}

// adds an equals condition
func (rf *RulerFilter) Eq(value interface{}) *RulerFilter {
	return rf.compare(eq, value)
}

// adds a not equals condition
func (rf *RulerFilter) Neq(value interface{}) *RulerFilter {
	return rf.compare(neq, value)
}

// adds a less than condition
func (rf *RulerFilter) Lt(value interface{}) *RulerFilter {
	return rf.compare(lt, value)
}

// adds a less than or equal condition
func (rf *RulerFilter) Lte(value interface{}) *RulerFilter {
	return rf.compare(lte, value)
}

// adds a greater than condition
func (rf *RulerFilter) Gt(value interface{}) *RulerFilter {
	return rf.compare(gt, value)
}

// adds a greater than or equal to condition
func (rf *RulerFilter) Gte(value interface{}) *RulerFilter {
	return rf.compare(gte, value)
}

// adds a matches (regex) condition
func (rf *RulerFilter) Matches(value interface{}) *RulerFilter {
	return rf.compare(matches, value)
}

// adds a greater than condition
func (rf *RulerFilter) NotMatches(value interface{}) *RulerFilter {
	return rf.compare(ncontains, value)
}

// comparator will either create a new ruler filter and add its filter
func (rf *RulerFilter) compare(comp int, value interface{}) *RulerFilter {
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
		rf = &RulerFilter{
			rf.Ruler,
			&Filter{
				comparator,
				rf.Path,
				value,
			},
		}
		// attach the new filter to the ruler
		rf.Ruler.filters = append(rf.Ruler.filters, rf.Filter)
	} else {
		//if there is no comparator, we can just set things on the current filter
		rf.Comparator = comparator
		rf.Value = value
	}

	return rf
}
