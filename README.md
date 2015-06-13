# go-ruler

json-y object testing in go

## Installation

```
go get github.com/hopkinsth/go-ruler
```

## Introduction

go-ruler is a reimplementaion of [ruler](https://github.com/RedVentures/ruler) in go, partly as an experiment, eventually as an optimization strategy. go-ruler does not yet support programmatically composing rules with function calls like ruler does and exists primarily for processing rules described arrays of JSON objects with this structure:

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