package pkg

func If2() {
	if true {
	}
}

func If7() { //@ diag("function has cyclomatic complexity of 7")
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
}

func ElseIf8() { //@ diag("function has cyclomatic complexity of 8")
	if true {
	} else if false {
	} else {
	}

	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
	if true {
	} else {
	}
}

func Switch8(i int) { //@ diag("function has cyclomatic complexity of 8")
	switch i {
	case 0:
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	}
}

func Compare6(i, j int) { //@ diag("function has cyclomatic complexity of 6")
	if i == j {
	} else if i+1 == j {
	} else if i+2 == j {
	} else if i+3 == j || i+4 == j {
	}
}
