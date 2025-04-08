package entity

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(id string) (ID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}
	return ID(uid), nil
}
