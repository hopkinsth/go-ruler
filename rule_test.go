package ruler

import "testing"

func TestPluck(t *testing.T) {
	exps := []struct {
		o       map[string]interface{}
		seeking string
		value   interface{}
	}{
		// test extracting a simple property
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.hello",
			"bob",
		},
		// test getting a nonexistent property
		// on an existing object
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.nope",
			nil,
		},
		//test getting a nonexistent property
		//for a nonexistent object
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.what.something.very.important",
			nil,
		},
		//test getting a property on a thing
		//that isn't or doesn't assert to be a map
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.hello.something.very.important",
			nil,
		},
		//test plucking something that isn't
		//a string
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"sup": 1,
				},
			},
			"hey.sup",
			1,
		},
		//test plucking where
		//the base obj doesn't exist
		{
			map[string]interface{}{},
			"hey.lol",
			nil,
		},
	}

	for _, e := range exps {
		res := pluck(e.o, e.seeking)
		if res != e.value {
			t.Error("failed to pluck!")
		}
	}
}
