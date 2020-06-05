package main

import (
	"database/sql"
	"github.com/AgileBits/go-redis-queue/redisqueue"
	"github.com/iseroukhov/brave-new-billing/pkg/entities/payment"
	"github.com/iseroukhov/brave-new-billing/pkg/server"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	s := server.New(cfg)
	mysql, err := s.MysqlDB()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := s.RedisQueue()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("queue is starting")
	for {
		job, err := queue.Pop()
		if err != nil {
			log.Println("queue pop error: ", err)
			time.Sleep(1 * time.Second)
		}
		if job != "" {
			switch {
			case strings.HasPrefix(job, "payment_timeout:"):
				id, err := strconv.ParseInt(strings.ReplaceAll(job, "payment_timeout:", ""), 10, 64)
				if err != nil {
					continue
				}
				if err := PaymentTimeout(id, mysql, queue); err != nil {
					log.Println("queue PaymentTimeout error: ", err)
				}
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func PaymentTimeout(pID int64, mysql *sql.DB, queue *redisqueue.Queue) error {
	paymentRepo := payment.NewRepository(mysql, queue)
	p, err := paymentRepo.GetByID(pID)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}
	_, err = paymentRepo.Cancell(p)
	return err
}
