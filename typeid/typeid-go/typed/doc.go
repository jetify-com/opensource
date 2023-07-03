// Package typed implements a statically checked version of TypeIDs
//
// To use it, define your own ID types that implement the IDType interface:
//
//	type UserIDType struct{}
//	func (UserIDType) Type() string { return "user" }
//
//	type AccountIDType struct{}
//	func (AccountIDType) Type() string { return "account" }
//
// And now you can use your IDTypes via generics. For example, to create a
// new ID of type user:
//
//	  import (
//		   typeid "go.jetpack.io/typeid/typed"
//		 )
//
//	  user_id, _ := typeid.New[UserIDType]()
//
// Because this implementation uses generics, the go compiler itself will
// enforce that you can't mix up your ID types. For example, a function with
// the signature:
//
//	func f(id typed.TypeID[UserIDType]) {}
//
// Will fail to compile if passed an id of type AccountIDType.
package typed
