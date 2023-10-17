package orm

import "github.com/zedisdog/brynn/orm/fields"

type User struct {
	TableName        string
	attributes       map[string]any
	originAttributes map[string]any
	fields           []fields.IField
	options          map[string]string
}

func (u *User) SetUsername(username string) {
	u.attributes["username"] = username
}

func (u *User) SetPassword(password string) {
	u.attributes["password"] = password
}

func (u *User) GetUsername() string {
	return u.attributes["username"].(string)
}

func (u *User) GetPassword() string {
	return u.attributes["password"].(string)
}

func (u User) dirty(field string) bool {
	attr, ok := u.attributes[field]
	oriAttr, oriOk := u.originAttributes[field]
	return ok != oriOk || attr != oriAttr
}

func (u User) DirtyMap() (m map[string]any) {
	m = make(map[string]any, len(u.attributes))
	for field, value := range u.attributes {
		if u.dirty(field) {
			m[field] = value
		}
	}

	return
}
