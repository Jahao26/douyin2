package repository

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type InteractiveEvent struct {
	Uid         int64  `json:"uid"`
	To_uid      int64  `json:"to_uid"`
	Action_type string `json:"action_type"`
}

var FollowWriter *kafka.Writer
var FollowConsumer *kafka.Reader
var FavoriteWriter *kafka.Writer
var FavoriteConsumer *kafka.Reader

func InitKafka() error {
	kafkaBrokers := "60.204.170.108:9092"
	FollowWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers),
		Topic:    "follow_events",
		Balancer: &kafka.Hash{},
	}
	FavoriteWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers),
		Topic:    "favorite_events",
		Balancer: &kafka.Hash{},
	}

	FollowConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"60.204.170.108:9092"},
		GroupID:  "follow-consumer-group",
		Topic:    "follow_events",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})
	FavoriteConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"60.204.170.108:9092"},
		GroupID:  "follow-consumer-group",
		Topic:    "favorite_events",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})

	go RunConsumer(FollowConsumer)
	go RunConsumer(FavoriteConsumer)
	fmt.Println("******************************************")
	fmt.Println("Kafka Init successfully")
	fmt.Println("******************************************")
	return nil
}

func RunConsumer(consumer *kafka.Reader) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			log.Println("Received termination signal, closing Kafka consumer")
			return
		default:
			m, err := consumer.ReadMessage(c)
			if err != nil {
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}

			var event InteractiveEvent
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				log.Printf("Error decoding message: %v\n", err)
				continue
			}

			err = updateMySQL(event)
			if err != nil {
				log.Printf("Error updating MySQL: %v\n", err)
			}
		}
	}
}

func updateMySQL(event InteractiveEvent) error {
	uid := event.Uid
	toUid := event.To_uid
	actionType := event.Action_type

	if strings.HasPrefix(actionType, "follow+") {
		action := strings.TrimPrefix(actionType, "follow+")

		if action == "1" {
			if err := NewRalationDao().AddRalation(uid, toUid); err != nil {
				return err
			}
			if err := NewUserDao().AddFollow(uid); err != nil {
				return err
			}
			if err := NewUserDao().AddFollower(toUid); err != nil {
				return err
			}
		} else {
			if err := NewRalationDao().RmRalation(uid, toUid); err != nil {
				return err
			}
			if err := NewUserDao().RmFollow(uid); err != nil {
				return err
			}
			if err := NewUserDao().RmFollower(toUid); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(actionType, "favorite+") {
		action := strings.TrimPrefix(actionType, "favorite+")

		if action == "1" { // 为了复用交互结构体，toUid=Vid
			if err := NewFavoriteDao().AddFavorite(uid, toUid); err != nil {
				return err
			}
			if err := NewUserDao().AddfavoriteCount(uid); err != nil {
				return err
			}
			if err := NewVideoDao().AddVideoFavorite(toUid); err != nil {
				return err
			}
		} else {
			if err := NewFavoriteDao().RmFavorite(uid, toUid); err != nil {
				return err
			}
			if err := NewUserDao().RmfavoriteCount(uid); err != nil {
				return err
			}
			if err := NewVideoDao().RmVideoFavorite(toUid); err != nil {
				return err
			}
		}
	}
	return nil
}
