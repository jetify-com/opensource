// Package typed implements a statically checked version of TypeIDs
//
// To use it, define your own ID types by creating a new string type:
//
// type UserID string
// const userPrefix = UserID("user")
//

// And now you can use your IDTypes with type enforcement. For example, to
// create a new ID of type user:
//
//	import (
//		typeid "go.jetpack.io/typeid/typed"
//	)
//
//	user_id, _ := typeid.T(userPrefix).New()
//
// Because this implementation uses generics, the go compiler itself will
// enforce that you can't mix up your ID types. For example, a function with
// the signature:
//
//	func f(id UserID) {}
//
// Will fail to compile if passed an id of type AccountID.
package typed
