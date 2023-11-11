package typeid_test

import (
	"fmt"

	"go.jetpack.io/typeid"
)

type UserID struct {
	typeid.TypeID
}

// UserID and other id subtypes need to implement typeid.Subtype
var _ typeid.Subtype = (*UserID)(nil)

// Use AllowedPrefix to define the prefix string for UserIDs
func (UserID) AllowedPrefix() string {
	return "user"
}

// Now do the same for AccountIDs

type AccountID struct {
	typeid.TypeID
}

var _ typeid.Subtype = (*AccountID)(nil)

func (AccountID) AllowedPrefix() string {
	return "account"
}

func Example() {
	userID, _ := typeid.New[UserID]()
	accountID, _ := typeid.New[AccountID]()

	// Each ID should have the correct type prefix:
	fmt.Printf("User ID prefix: %s\n", userID.Prefix())
	fmt.Printf("Account ID prefix: %s\n", accountID.Prefix())

	// The compiler considers their go types to be different:
	fmt.Printf("%T != %T\n", userID, accountID)
	// Output:
	// User ID prefix: user
	// Account ID prefix: account
	// typeid_test.UserID != typeid_test.AccountID
}
