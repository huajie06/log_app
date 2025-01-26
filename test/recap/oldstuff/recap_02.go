package recap

import "fmt"

type person struct {
	name string
	age  int
}

func (p person) String() string {
	return fmt.Sprintf("Hi! the perosnal's name is %v, age is %v", p.name, p.age)
}

func Recap02() {
	a, b := learnMemory()
	fmt.Println(a, b)

	x := 100

	bigB := func() bool {
		return x > 9
	}

	fmt.Println("My func has value:", bigB())

	fmt.Println("-----------------------------")
	pp := person{"alex", 100}
	fmt.Println(pp.String())
}

// Go is fully garbage collected. It has pointers but no pointer arithmetic.
// You can make a mistake with a nil pointer, but not by incrementing a pointer.
// Unlike in C/Cpp taking and returning an address of a local variable is also safe.
func learnMemory() (p, q *int) {
	// Named return values p and q have type pointer to int.
	p = new(int) // Built-in function new allocates memory.
	// The allocated int slice is initialized to 0, p is no longer nil.

	s := make([]int, 20) // Allocate 20 ints as a single block of memory.
	s[3] = 7             // Assign one of them.
	r := -2              // Declare another local variable.
	return &s[3], &r     // & takes the address of an object.
}
