package pkg

func Test1() {
	const mName = "Test2" //@ diag(`const mName should use function name`)
	_ = mName
}
