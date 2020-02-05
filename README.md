# go-jnode

[![GoDoc](https://godoc.org/github.com/soluble-ai/go-jnode?status.svg)](https://pkg.go.dev/github.com/soluble-ai/go-jnode?tab=doc)
[![GoReport](https://goreportcard.com/badge/github.com/soluble-ai/go-jnode)](https://goreportcard.com/report/github.com/soluble-ai/go-jnode)

`go-jnode` is a Go module that makes using generic JSON structures easy.
Its API is modeled on Java's [jackson tree model](https://github.com/FasterXML/jackson-databind).

Go's builtin generic JSON object uses `map[string]interface{}` for objects,
`[]interface{}` for arrays, and `string`, `bool`, `float64` for everything else.  However, using these generic interfaces safely is awkward and verbose.

`go-jnode` makes things a little easier:

* Create a `*jnode.Node` with any of the factory methods e.g. `jnode.NewObjectNode()` or `jnode.FromJSON()`
* Use chained `n.Path(field)` or `n.Get(index)` calls to navigate an object.
* Use `n.Entries()` to iterate over maps, and `n.Elements()` to iterate over arrays.
* Use `n.AsText()` to get a text value. (Or `n.AsBool()`, `n.AsInt()` etc)
* Navigation is safe - if the object doesn't have a field or an array doesn't have an index a single `MissingNode` is returned, for which `n.IsMissing()` returns `true`.  (The text value of a missing node is empty.)

To install:

    go get github.com/soluble-ai/go-jnode

Short example:

```go
import "github.com/soluble-ai/go-jnode"

func example() {
  n := jnode.NewObjectNode().Put("greeting", "hello").Put("subject", "world")

  fmt.Printf("%v\n", n)                     // {"greeting":"hello","subject":"world"}
  fmt.Println(n.Path("greeting").AsText())  // greeting
  fmt.Println(n.Path("not-there").AsText()) // <empty-string>

  n.PutArray("list").Append(1).Append(2)    // adds an Array
  fmt.Println(n.Path("list").Get(1))        // 2

  // can also add slices, which get flattened
  n.Path("list").Append([]string{ "hello", "world" })

  // iteration: use Elements() for arrays, Entries() for maps
  for _, e := range n.Path("list").Elements() {
    fmt.Println(e.AsText());    // "hello" followed by "world"
  }

  // construct from a string
  m, _ := jnode.FromJSON(`{"code":200,"message":"howdy"})`)
  fmt.Println(m.Path("code").AsInt())       // 200
}
```
