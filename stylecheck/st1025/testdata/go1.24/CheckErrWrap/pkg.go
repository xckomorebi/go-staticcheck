package pkg

import "example.com/xcerr"

const MName = "MyMethod"

type MyStruct struct {
}

const nMystruct = "MyStruct"

func (m MyStruct) MyMethod() error {
	const mName = "MyMethod"

	return xcerr.New(nMystruct, mName, "error")
}
