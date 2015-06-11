package ruler

import "testing"

func TestPluck(t *testing.T) {
	o := make(map[string]interface{})
	o["hey"] = make(map[string]interface{})
	o["hey"].(map[string]interface{})["hello"] = "bob"

	r := pluck(o, "hey.hello")
	if r != "bob" {
		t.Error("didn't pluck properly!")
	}
}
