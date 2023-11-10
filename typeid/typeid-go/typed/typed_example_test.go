package typed_test

import (
	"fmt"
	"time"

	untyped "go.jetpack.io/typeid"
	typeid "go.jetpack.io/typeid/typed"
)

type userPrefix struct{}

func (userPrefix) Type() string { return "user" }

type UserID struct{ typeid.TypeID[userPrefix] }

type accountPrefix struct{}

func (accountPrefix) Type() string { return "account" }

type AccountID struct{ typeid.TypeID[accountPrefix] }

type OrgID struct {
	untyped.TypeID `prefix:"org"`
}

func Example() {
	userID, _ := typeid.New[UserID]()
	accountID, _ := typeid.New[AccountID]()
	orgID, _ := typeid.New[OrgID]()
	orgID2, _ := untyped.New2[OrgID]()
	// Each ID should have the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", userID.Type())
	fmt.Printf("Account ID prefix: %s\n", accountID.Type())
	fmt.Printf("Org ID prefix: %s\n", orgID.Type())
	fmt.Printf("Org2 ID prefix: %s\n", orgID2.Type())
	// The compiler considers their go types to be different:
	fmt.Printf("%T != %T\n", userID, accountID)

	start := time.Now()
	for i := 0; i < 1000000; i++ {
		id, _ := typeid.New[UserID]()
		_ = id.Type()
	}
	fmt.Printf("1000000 New[UserID] calls took %v\n", time.Since(start))

	start = time.Now()
	for i := 0; i < 1000000; i++ {
		id, _ := typeid.New[OrgID]()
		_ = id.Type()
	}
	fmt.Printf("1000000 New[OrgID] calls took %v\n", time.Since(start))

	start = time.Now()
	for i := 0; i < 1000000; i++ {
		id, _ := untyped.New2[OrgID]()
		_ = id.Type()
	}
	fmt.Printf("1000000 New2[OrgID] calls took %v\n", time.Since(start))

	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typed.TypeID[go.jetpack.io/typeid/typed_test.UserID] != typed.TypeID[go.jetpack.io/typeid/typed_test.AccountID]
}
