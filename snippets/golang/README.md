# Go

- [Switch Case Behavior](#switch-case-behavior)
- [Compile-Time Checks](#compile-time-checks)
- [Arrays vs Slices](#arrays-vs-slices)
- [Nil Slice vs Empty Slice in JSON](#nil-slice-vs-empty-slice-in-json)

## Switch Case Behavior

- Go's switch cases **break automatically by default** - no `break` statement needed
- Use `fallthrough` keyword to explicitly continue to the next case
- The `default` case **always executes last**, regardless of its position in the code

```go
switch value {
case 1:
    fmt.Println("One")
    // Automatically breaks
case 2:
    fmt.Println("Two")
    fallthrough  // Explicitly continue to next case
case 3:
    fmt.Println("Three")
default:
    fmt.Println("Default")  // Runs last even if written first
}
```

## Compile-Time Checks

- Go **detects circular dependencies at compile time**
- Prevents import cycles between packages
- Compilation fails immediately with error: `"import cycle not allowed"`

```go
// package A imports package B
// package B imports package A
// → Compile error: "import cycle not allowed"
```

## Arrays vs Slices

### Arrays
- **Fixed-size**, **value type**
- Assignment **creates a full copy**
- Changes to copy don't affect original

```go
arr1 := [3]int{1, 2, 3}
arr2 := arr1           // Full copy
arr2[0] = 100

fmt.Println(arr1)      // [1 2 3]    (unchanged)
fmt.Println(arr2)      // [100 2 3]
```

### Slices
- **Dynamic-size**, **reference type**
- Assignment **shares the underlying array**
- Changes through one variable affect all references

```go
slice1 := []int{1, 2, 3}
slice2 := slice1         // Both point to same data
slice2[0] = 100

fmt.Println(slice1)      // [100 2 3]  (changed!)
fmt.Println(slice2)      // [100 2 3]
```

### Creating Independent Slice Copies

Use `copy()` to create an independent slice:

```go
slice1 := []int{1, 2, 3}
slice2 := make([]int, len(slice1))
copy(slice2, slice1)    // Deep copy

slice2[0] = 100
fmt.Println(slice1)     // [1 2 3]    (unchanged)
```

## Nil Slice vs Empty Slice in JSON

- **Nil slice** marshals to `null`
- **Empty slice** marshals to `[]`
- When unmarshaling, `null` creates a nil slice, `[]` creates an empty slice

```go
type Data struct {
    Items []string `json:"items"`
}

// Nil slice → null
d1 := Data{Items: nil}
json1, _ := json.Marshal(d1)
fmt.Println(string(json1))  // {"items":null}

// Empty slice → []
d2 := Data{Items: []string{}}
json2, _ := json.Marshal(d2)
fmt.Println(string(json2))  // {"items":[]}

// Unmarshaling null → nil slice
var d3 Data
json.Unmarshal([]byte(`{"items":null}`), &d3)
fmt.Println(d3.Items == nil)  // true

// Unmarshaling [] → empty slice
var d4 Data
json.Unmarshal([]byte(`{"items":[]}`), &d4)
fmt.Println(d4.Items == nil)  // false
fmt.Println(len(d4.Items))    // 0
```