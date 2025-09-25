package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dangLuan01/user-manager/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQService struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

func NewRabitMQService(amqpURL string) (RabbitMQService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Println(err)
		return nil, utils.WrapError(string(utils.ErrCodeInternal), "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Println(err)
		return nil, utils.WrapError(string(utils.ErrCodeInternal), "Failed to open channel", err)
	}

	return &rabbitMQService{
		conn: conn,
		channel: ch,
	}, nil
}

func (r *rabbitMQService) Publish(ctx context.Context, queue string, message any) error {
	_, err := r.channel.QueueDeclare(queue, true, false, false, false, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return err
	}
	err = r.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *rabbitMQService) Consume(ctx context.Context, queue string, handler func ([]byte) error) error{

	_, err := r.channel.QueueDeclare(queue, true, false, false, false, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	msgs, err := r.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	
	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return 
				}
				if err := handler(msg.Body); err != nil {
					msg.Nack(false, false)
				} else {
					msg.Ack(false)
				}
			case <- ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *rabbitMQService) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}