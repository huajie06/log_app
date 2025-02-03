package recap

import "fmt"

func Recap01() {
	fmt.Println("hello world")

	fmt.Println(doStuff(100, 200))

	result := doStuff(9, 10)
	fmt.Println(result)

	p1, p2, p3 := doStuffThree(1, 2)
	fmt.Println(p1, p2, p3)

	str := "hello world"
	str02 := `what the heck??
why i can not have a line?`

	fmt.Println(str, str02)

	arr := [4]int{0, 0, 0, 3}
	fmt.Println(arr)

	arr01 := arr
	arr01[0] = 1
	fmt.Println(arr01)

	// slice, with dynamic size
	s3 := []int{4, 5, 9}
	fmt.Println(s3)

	fmt.Println("--------------------")

	bs := []byte("a slice")
	fmt.Println(bs)

	fmt.Println("--------------------")
	fmt.Println([]byte("a"))
	fmt.Println('a')

	fmt.Println("--------------------")
	s := []int{1, 2, 3}
	s = append(s, 4, 5, 6)
	fmt.Println(s)

	s = append(s, []int{7, 8, 9}...)
	fmt.Println(s)
	fmt.Println("--------------------")

	for i := 1; i < 5; i++ {
		fmt.Println(i, "what!!!")
	}

	for k, v := range map[string]int{"a": 1, "b": 2, "c": 3} {
		fmt.Printf("printing dict -- string: %s, value: %v.\n", k, v)
	}

	fmt.Println("------------------")
	for i, v := range []string{"a", "b"} {
		fmt.Printf("printing slice -- index: %v, value: %v\n", i, v)
	}
}

func doStuff(x, y int) int {
	return x + y
}

func doStuffThree(x, y int) (sum, prod, chg int) {
	return x + y, x * y, x + x + y*y

}
