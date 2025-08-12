package database

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:public.users,alias:u"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	FirstName     string    `bun:"first_name,notnull,type:text"`
	Surname       string    `bun:"surname,notnull,type:text"`
	Patronymic    *string   `bun:"patronymic,type:text"`
	Age           int       `bun:"age,notnull,type:integer"`
	Nationality   string    `bun:"nationality,notnull,type:text"`
	Gender        string    `bun:"gender,notnull,type:text"`
	Emails        []string  `bun:"emails,notnull,type:jsonb"`
	//Friends       []*User   `bun:"m2m:friends,join:User=Friend"`
}

type Friend struct {
	bun.BaseModel `bun:"table:public.friends,alias:f"`
	PersonID      uuid.UUID `bun:"person_id,pk,type:uuid"`
	Person        *User     `bun:"rel:belongs-to,join:person_id=id"`
	FriendID      uuid.UUID `bun:"friend_id,pk,type:uuid"`
	Friend        *User     `bun:"rel:belongs-to,join:friend_id=id"`
}
