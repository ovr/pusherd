package main

type PushNotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

type PushNotificationTask struct {
	UserId  uint64                  `json:"userId"`
	Payload PushNotificationPayload `json:"payload"`
}

type PushToken struct {
	Id  string `gorm:"primary_key"`
	Uid uint64
}
