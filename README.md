## Go Generics

- Go Generics are defined using type parameters and can be used on functions, structs, interfaces, and methods.
- Type parameters are specified within square brackets `[]` after the function or type name.
- Generic types must be constrained using type constraints to specify the allowed types for the type parameters.
	- **Constrained by interface**: Define an interface that needs to be implemented by the type parameter.
		- Example:
		```go
		type Stringer interface {
			String() string
		}

		func PrintString[T Stringer](value T) {
			fmt.Println(value.String())
		}
		```

		> The special `comparable` constraint allows any type that supports comparison operators (`==`, `!=`). This is useful when using the generic type as a key in maps or when checking for equality.
		- Example:
			```go
			func AreEqual[T comparable](a, b T) bool {
				return a == b
			}
			```
	- **Constrained by type**: Define a (base) type that the type parameter must be based on.
		- Example:
		```go
		func PrintValue[T ~int](value T) {
			fmt.Println(value)
		}

		type MyInt int

		PrintValue(MyInt(42)) // Valid, MyInt's underlying type is int
		PrintValue(100)       // Valid, int is allowed
		```
		> The `~` operator allows the type parameter to accept any type whose underlying type is the specified base type.

	- **Constrained by multiple types (type sets)**: Define a set of types that a type parameter can accept.
		- Example: 
		```go
		func PrintType[T int | string](value T) {
			fmt.Println(value)
		}

		//or

		func SumNumbers[T interface {int | float64}](a, b T) T {
			return a + b
		}

		PrintType(42)        // Valid, int is allowed
		PrintType("Hello")   // Valid, string is allowed
		```
		> The Go experimental library `golang.org/x/exp/constraints` provides some predefined type sets like `constraints.Ordered` for types that support ordering operators (`<`, `>`, etc.) or `constraints.Integer` for integer types.

- Go can infer type parameters when calling generic functions, so you often don't need to specify them explicitly.
	- Example:
	```go
	func PrintValue[T any](value T) {
		fmt.Println(value)
	}		
	PrintValue(42)          // Type parameter T is inferred as int
	PrintValue("Hello")     // Type parameter T is inferred as string
	```
- At compile time, Go uses **Monomorphization** to generate type specific versions of generic functions/types for each unique set of type arguments used in the code. This ensures that there is no runtime overhead associated with using generics. This will however lead to larger binary sizes if many different type arguments are used.

## Go Generics Limitations
- Methods cannot have their own type parameters; only the type they are defined on (receiver) can have type parameters.
	```go
	func (t T) MethodName[U any](param U) { 
		// This is not allowed 
	}
	```
- Methods cannot constrain their receiver type parameters further.
	```go
	type Container[T any] struct {
		value T
	}

	func (c Container[T comparable]) IsEqual(other Container[T]) bool { 
		// This is not allowed 
		return c.value == other.value
	}
	```
- You cannot implement the `comparable` constraint on custom types.
- Limited type assertion, even when using `any` as a constraint.
	```go
	func ProcessValue[T any](value T) {
		str,ok := value.(string) // This is not allowed
	}

	// Workaround:
	func ProcessValue[T any](value T) {
		str, ok := any(value).(string) // This is allowed
	}
	```
- Methods or fields of struct-constraints cannot be accessed on the type parameter directly.
	```go
	type Box struct {
		value int
	}

	func (b Box) GetValue() int {
		return b.value
	}

	func ProcessBox[T Box](box T) {
		val := box.GetValue() // This is not allowed
		val := box.value // This is not allowed
	}

	// Workaround:
	func ProcessBox[T Box](box T) {
		val := Box(box).GetValue() // This is allowed
	}
	```
- Multiple interfaces cannot be used as a type set constraint.
	```go
	type Reader interface {
		Read(p []byte) (n int, err error)
	}

	type Writer interface {
		Write(p []byte) (n int, err error)
	}

	func ReadWrite[T Reader | Writer](rw T) { 
		// This is not allowed 
	}
	```

