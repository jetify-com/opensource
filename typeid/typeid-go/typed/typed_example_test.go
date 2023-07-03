package typed_test

import (
	"fmt"

	typeid "go.jetpack.io/typeid/typed"
)

type UserIDType struct{}

func (UserIDType) Type() string { return "user" }

type AccountIDType struct{}

func (AccountIDType) Type() string { return "account" }

func Example() {
	user_id, _ := typeid.New[UserIDType]()
	account_id, _ := typeid.New[AccountIDType]()
	// Each ID should have the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", user_id.Type())
	fmt.Printf("Account ID prefix: %s\n", account_id.Type())
	// The compiler considers their go types to be different:
	fmt.Printf("%T != %T\n", user_id, account_id)

	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typed.TypeID[go.jetpack.io/typeid/typed_test.UserIDType] != typed.TypeID[go.jetpack.io/typeid/typed_test.AccountIDType]
}
