package typeid

import (
	"github.com/gofrs/uuid"
	"go.jetpack.io/typeid/base32"
)

type TypeID struct {
	Type string
	UUID [16]byte
}

var Nil = TypeID{
	UUID: uuid.Nil,
}

func New(prefix string) (TypeID, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return Nil, err
	}
	tid := TypeID{
		Type: prefix,
		UUID: uid,
	}
	return tid, nil
}

func (tid TypeID) UUIDString() string {
	return base32.Encode(tid.UUID)
}

func (tid TypeID) String() string {
	return tid.Type + "_" + tid.UUIDString()
}

func Must(tid TypeID, err error) TypeID {
	if err != nil {
		panic(err)
	}
	return tid
}
