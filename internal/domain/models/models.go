package models

type Users struct {
	TgId      uint64 `json:"tg_id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	LastName  string `json:"last_name"`
}

type Topics struct {
	Topic string `json:"topic_name"`
	TgId  uint64 `json:"tg_id"`
}

type Words struct {
	Word      string `json:"word"`
	Translate string `json:"translate"`
	TopicName string `json:"topicName"`
}
