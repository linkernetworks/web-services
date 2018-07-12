package entity

import "gopkg.in/mgo.v2/bson"

const UserCollectionName string = "users"

type User struct {
	ID                    bson.ObjectId `bson:"_id" json:"id"`
	SerialNumber          string        `bson:"serial_number" json:"serial_number"`
	Email                 string        `bson:"email,omitempty" json:"email"`
	Password              string        `bson:"password,omitempty" json:"password,omitempty"`
	FirstName             string        `bson:"first_name" json:"first_name"`
	LastName              string        `bson:"last_name" json:"last_name"`
	CountryCode           string        `bson:"country_code" json:"country_code"`
	Cellphone             string        `bson:"cellphone" json:"cellphone"`
	Roles                 []string      `bson:"roles" json:"roles"`
	Verified              bool          `bson:"verified" json:"verified"`
	VerificationCode      string        `bson:"verification_code" json:"verification_code"`
	Jwt                   string        `bson:"jwt" json:"jwt"`
	AccessToken           string        `bson:"access_token" json:"access_token"`
	AccessTokenExpiryTime int64         `bson:"access_token_expiry_time" json:"access_token_expiry_time"`
	RefreshToken          string        `bson:"refresh_token" json:"refresh_token"`
	CreatedAt             int64         `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt             int64         `bson:"updated_at,omitempty" json:"updated_at"`
	LastLoggedInAt        int64         `bson:"last_loggedin_at,omitempty" json:"last_loggedin_at"`
	Revoked               bool          `bson:"revoked" json:"revoked"`
	Preference            Preference    `bson:"preference" json:"preference"`
}

func (u *User) GetCollection() string {
	return UserCollectionName
}
