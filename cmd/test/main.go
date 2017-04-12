package main

import (
	"gopkg.in/maddevsio/fcm.v1"
	"log"
	"fmt"
)

func main() {
	data := map[string]string{}

	c := fcm.NewFCM("")
	token := ""

	for {
		response, err := c.Send(&fcm.Message{
			Data:             data,
			RegistrationIDs:  []string{token},
			ContentAvailable: true,
			Priority:         fcm.PriorityHigh,
			Notification: &fcm.Notification{
				Title: "Hello",
				Body:  "Antonio!",
				Sound: "default",
			},
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Status Code   :", response.StatusCode)
		fmt.Println("Success       :", response.Success)
		fmt.Println("Fail          :", response.Fail)
		fmt.Println("Canonical_ids :", response.CanonicalIDs)
		fmt.Println("Topic MsgId   :", response.MsgID)
		fmt.Println("Topic Results   :", response.Results)
	}
}
