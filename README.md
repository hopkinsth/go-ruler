# go-ruler

json-y object testing in go

## Installation

```
go get github.com/hopkinsth/go-ruler
```

## Introduction

go-ruler is an implementaion of [ruler](https://github.com/RedVentures/ruler) in go, partly as an experiment, eventually as an optimization strategy. go-ruler supports programmatically constructing rules in a way similar to js-ruler and can also process rules stored in arrays of JSON objects with this structure:

```json
{
  "comparator": "eq",
  "path": "library.name",
  "value": "go-ruler"
}
```

Each of those objects (or 'filters') describes a condition for a property on some JSON object or json-object-esque structure. (In this go implementation, we're using `map[string]interface{}`.)

## Example

```go
package main

import "github.com/hopkinsth/go-ruler"
import "fmt"

func main() {
  rules := []byte(`[
    {"comparator": "eq", "path": "library.name", "value": "go-ruler"},
    {"comparator": "gt", "path": "library.age", "value": 0.5}
  ]`)

  engine, _ := ruler.NewRulerWithJson(rules)

  result := engine.Test(map[string]interface{}{
    "library": map[string]interface{}{
      "name": "go-ruler",
      "age":  1.24,
    },
  })

  fmt.Println(result == true)
}
```

## License
Copyright 2015 Thomas Hopkins

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this work except in compliance with the License. You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.