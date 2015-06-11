package ruler

import "testing"

func TestPluck(t *testing.T) {
	exps := []struct {
		o       map[string]interface{}
		seeking string
		value   interface{}
		name    string
	}{

		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.hello",
			"bob",
			`test extracting a simple property`,
		},

		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.nope",
			nil,
			`test getting a nonexistent property
			on an existing object`,
		},
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.what.something.very.important",
			nil,
			`test getting a nonexistent property
			for a nonexistent object`,
		},
		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"hello": "bob",
				},
			},
			"hey.hello.something.very.important",
			nil,
			`test getting a property on a thing
			that isn't or doesn't assert to be a map`,
		},

		{
			map[string]interface{}{
				"hey": map[string]interface{}{
					"sup": 1,
				},
			},
			"hey.sup",
			1,
			`test plucking something that isn't
			a string`,
		},
		{
			map[string]interface{}{},
			"hey.lol",
			nil,
			`test plucking where the base obj doesn't exist`,
		},
		{
			map[string]interface{}{
				"test": map[string]interface{}{
					"thing": map[string]interface{}{
						"here": map[string]interface{}{
							"today": map[string]interface{}{
								"is": map[string]interface{}{
									"awesome": map[string]interface{}{
										"with": map[string]interface{}{
											"thestuff": "no dice",
										},
									},
								},
							},
						},
					},
				},
			},
			"test.thing.here.today.is.awesome.with.thestuff",
			"no dice",
			"testing deeply nested property",
		},
	}

	for _, e := range exps {
		res := pluck(e.o, e.seeking)
		if res != e.value {
			t.Errorf(
				"error with: %s\nfailed to pluck! %s != %s",
				e.name,
				res,
				e.value,
			)
		}
	}
}

func BenchmarkPluckShallow(b *testing.B) {
	o := map[string]interface{}{
		"hey": map[string]interface{}{
			"there": 4,
		},
	}

	for i := 0; i < b.N; i += 1 {
		r := pluck(o, "hey.there")
		if r != 4 {
			b.Errorf("fail bench, val was %s", r)
		}
	}
}

func BenchmarkPluckDeep(b *testing.B) {
	o := map[string]interface{}{
		"test": map[string]interface{}{
			"thing": map[string]interface{}{
				"here": map[string]interface{}{
					"today": map[string]interface{}{
						"is": map[string]interface{}{
							"awesome": map[string]interface{}{
								"with": map[string]interface{}{
									"thestuff": "no dice",
								},
							},
						},
					},
				},
			},
		},
	}

	for i := 0; i < b.N; i += 1 {
		r := pluck(o, "test.thing.here.today.is.awesome.with.thestuff")
		if r != "no dice" {
			b.Errorf("fail bench, val was %s", r)
		}
	}
}
