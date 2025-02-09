package db_models

type Users struct {
	TgId      int64  `json:"tg_id" db:"tg_id"`
	Username  string `json:"username" db:"username"`
	Firstname string `json:"firstname" db:"firstname"`
	LastName  string `json:"last_name" db:"last_name"`
}

type Topics struct {
	Topic string `json:"topic_name" db:"topic_name"`
	TgId  int64  `json:"tg_id" db:"tg_id"`
}

type Words struct {
	Word      string `json:"word" db:"word"`
	Translate string `json:"translate" db:"translate"`
	TopicName string `json:"topicName" db:"topicName"`
}
