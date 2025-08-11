package database

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID          uuid.UUID `bun:",pk"`
	FirstName   string
	Surname     string
	Patronymic  *string
	Gender      string
	Nationality string
	Age         int
	Email       []*Email `bun:"rel:has-many,join:id=user_id"`
	Friends     []*User  `bun:"m2m:friends,"`
}

type Email struct {
	bun.BaseModel

	UserID  uuid.UUID
	Address string
}

type Friend struct {
	bun.BaseModel

	UserID   uuid.UUID `bun:"user_id,pk"`
	User     *User     `bun:"rel:belongs-to,join:user_id=id"`
	FriendID uuid.UUID `bun:"friend_id,pk"`
	Friend   *User     `bun:"rel:belongs-to,join:friend_id=id"`
}
