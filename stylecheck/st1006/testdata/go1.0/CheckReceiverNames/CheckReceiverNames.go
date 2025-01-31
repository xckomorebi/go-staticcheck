// Package pkg ...
package pkg

type T1 int

func (x T1) Fn1()    {} //@ diag(`receiver name should use the first letter of its type`)
func (y T1) Fn2()    {} //@ diag(`receiver name should use the first letter of its type`)
func (x T1) Fn3()    {} //@ diag(`receiver name should use the first letter of its type`)
func (T1) Fn4()      {}
func (_ T1) Fn5()    {} //@ diag(`receiver name should not be an underscore, omit the name if it is unused`)
func (self T1) Fn6() {} //@ diag(`receiver name should be a reflection of its identity`)
