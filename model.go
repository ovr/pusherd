package main

type SerializedJSONField struct {
	Scalar string      `json:"@scalar"`
	Value  interface{} `json:"@value"`
}

type SerializedJSONArray struct {
	Map   string                         `json:"@map"`
	Value map[string]SerializedJSONField `json:"@value"`
}

type SerializedPushNotificationTask struct {
	Type    string              `json:"@type"`
	UserId  SerializedJSONField `json:"userId"`
	Payload SerializedJSONArray `json:"payload"`
}

type PushToken struct {
	Id  string `gorm:"primary_key"`
	Uid uint64
}
