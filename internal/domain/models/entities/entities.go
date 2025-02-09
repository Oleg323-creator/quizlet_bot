package entities

type Users struct {
	TgId      uint64 `json:"tg_id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	LastName  string `json:"last_name"`
}
