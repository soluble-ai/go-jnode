# go-jnode

`go-jnode` is a Go module that makes using generic JSON structures easy.
Its API is modeled on Java's [jackson tree model](https://github.com/FasterXML/jackson-databind).

Go can unmarshal JSON into a generic `interface{}` using `map[string]interface{}` for objects,
`[]interface{}` for arrays, and `string`, `bool`, `float64` for everything else.
(And it back marshal it back of course.)  But manipulating the resulting values is awkard.

And building generic JSON objects in straight go is also awkard.  And verbose.

`go-jnode` makes both these things easy.

To use:

    import "github.com/soluble-ai/go-jnode"

Short example:

```go
n := jnode.NewObjectNode().Put("greeting", "hello").Put("subject", "world")

fmt.Printf("%v\n", n)                     // {"greeting":"hello","subject":"world"}
fmt.Println(n.Path("greeting").AsText())  // greeting
fmt.Println(n.Path("not-there").AsText()) // <empty-string>

n.PutArray("list").Append(1).Append(2)    // adds an Array
fmt.Println(n.Path("list").Get(1))        // 2

m, _ := jnode.FromJSON(`{"code":200,"message":"howdy"})`)

fmt.Println(m.Path("code").AsInt())       // 200
```
