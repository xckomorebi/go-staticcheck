package pkg

import "fmt"

func getNum() int {
	return 1
}

func f() {
	num := getNum() //@ diag("assigned too early")

	if true {
	}

	fmt.Println(num)
}

func f2() {
	num := getNum() //@ diag("assigned outside of usage scope")

	if true {
		fmt.Println(num)
	}
}

func f3() {
	num := getNum()

	if true {
		fmt.Println(num)
	}
	fmt.Println(num)
}

func f4() {
	num := getNum()

	if num2 := 4; true {
		fmt.Println(num2)
	}
	fmt.Println(num)
}
