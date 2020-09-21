package database

type User struct {
	Login string `json:"login"`
	HashPassword []byte `json:"password"`
	JWT string `json:"-"`
}
