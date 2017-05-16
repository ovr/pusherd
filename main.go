package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"gopkg.in/maddevsio/fcm.v1"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initDB(configuration *Configuration) *gorm.DB {
	db, err := gorm.Open(configuration.DB.Dialect, configuration.DB.Uri)
	if err != nil {
		panic(err)
	}

	db.LogMode(configuration.DB.ShowLog)
	db.DB().SetMaxIdleConns(configuration.DB.MaxIdleConnections)
	db.DB().SetMaxOpenConns(configuration.DB.MaxOpenConnections)

	return db
}

func initAMQP(configuration *Configuration) *amqp.Channel {
	conn, err := amqp.Dial(configuration.AMQP.Uri)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")

	err = ch.ExchangeDeclare(
		"interpals",
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare exchange")

	queue, err := ch.QueueDeclare(
		"push-notifications",
		true,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	err = ch.QueueBind(
		queue.Name,
		"push-notifications",
		"interpals",
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue")

	return ch
}

func pushSender(client *fcm.FCM, db *gorm.DB, receive chan SerializedPushNotificationTask) {
	for {
		for task := range receive {
			pushTokens := []PushToken{}

			db.Where("uid = ?", uint64(task.UserId.Value.(float64))).Find(&pushTokens)

			if len(pushTokens) == 0 {
				continue
			}

			registrationIDs := []string{}

			for _, token := range pushTokens {
				registrationIDs = append(registrationIDs, token.Id)
			}

			log.Print(registrationIDs)

			title := task.Payload.Value["title"].Value
			if title == nil {
				log.Println("%s: %s", "Title is nill", title)

				continue;
			}

			body := task.Payload.Value["body"].Value
			if body == nil {
				log.Println("%s: %s", "Body is nill", body)

				continue;
			}

			sound := task.Payload.Value["sound"].Value
			if sound == nil {
				log.Println("%s: %s", "Sound is nill", body)

				continue;
			}

			response, err := client.Send(&fcm.Message{
				Data:             nil,
				RegistrationIDs:  registrationIDs,
				ContentAvailable: true,
				Priority:         fcm.PriorityHigh,
				Notification: &fcm.Notification{
					Title: title.(string),
					Body:  body.(string),
					Sound: sound.(string),
				},
			})

			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Status Code   :", response)
				fmt.Println("Status Code   :", response.StatusCode)
				fmt.Println("Success       :", response.Success)
				fmt.Println("Fail          :", response.Fail)
				fmt.Println("Canonical_ids :", response.CanonicalIDs)
				fmt.Println("Topic MsgId   :", response.MsgID)
				fmt.Println("Topic Results   :", response.Results)
			}
		}
	}
}

func consume(ch *amqp.Channel, tasksChannel chan SerializedPushNotificationTask) {
	receiveChannel, err := ch.Consume(
		"push-notifications", // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	failOnError(err, "Failed to register a consumer")

	for {
		for d := range receiveChannel {
			log.Printf("Received a message: %s \n", d.Body)

			task := SerializedPushNotificationTask{}

			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				d.Nack(false, false)

				log.Println(err)
			} else {
				tasksChannel <- task

				d.Ack(false)
			}
		}
	}
}

func main() {
	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "./config.json", "Config filepath")
	flag.Parse()

	configuration := &Configuration{}
	configuration.Init(configFile)

	client := fcm.NewFCM(configuration.FCM.Key)
	db := initDB(configuration)
	ch := initAMQP(configuration)

	tasksChannel := make(chan SerializedPushNotificationTask, configuration.Buffer)

	for i := 0; i < configuration.Senders; i++ {
		go pushSender(client, db, tasksChannel)
	}

	consume(ch, tasksChannel)
}
