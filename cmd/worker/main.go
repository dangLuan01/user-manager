package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dangLuan01/user-manager/internal/app"
	"github.com/dangLuan01/user-manager/internal/config"
	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/dangLuan01/user-manager/pkg/mail"
	"github.com/dangLuan01/user-manager/pkg/rabbitmq"
)

type Worker struct {
	rabbitMQ rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	cfg *config.Config
}

func NewWorker(cfg *config.Config) *Worker {
	rabbitMQ, err := rabbitmq.NewRabitMQService(utils.GetEnv("RABBITMQ_URL","amqp://guest:guest@localhost:5672"))
	if err != nil {
		log.Printf("Failed to init rabbitMQ service:%s", err)
	}

	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		log.Fatalf("⛔ Unable to init mail:%s", err)
	}

	mailService, err := mail.NewMailService(cfg, factory)
	if err != nil {
		log.Fatalf("⛔ Unable to init mail service:%s", err)
	}

	return &Worker{
		rabbitMQ: rabbitMQ,
		mailService: mailService,
		cfg: cfg,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	const emailQueueName = "auth_email_queue"

	hanldler := func(body []byte) error {

		var email mail.Email

		if err := json.Unmarshal(body, &email); err != nil {
			log.Printf("Failed to unmarshal message:%s", err)
			return err
		}

		if err := w.mailService.SendMail(ctx, &email); err != nil {
			log.Println(err)
			return utils.NewError(string(utils.ErrCodeInternal), "Failed to send email.")
		}

		return nil
	}

	if err := w.rabbitMQ.Consume(ctx, emailQueueName, hanldler); err != nil {
		log.Printf("Failed to start consumer:%s", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (w *Worker) Shutdown(ctx context.Context) error {
	if err := w.rabbitMQ.Close(); err != nil {
		log.Printf("Failed to close rabbitMq:%s", err)
		return err
	}

	select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Println("Shutdown timeout exceeded")
				return ctx.Err()
			}
		default:
	}

	return nil
}

func main() {
	
	app.LoadEnv()
	
	cfg := config.NewConfig()

	worker := NewWorker(cfg)
	if worker == nil {
		log.Println("Failed to create worker")
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := worker.Start(ctx); err != nil && err != context.Canceled {
			log.Println("Worker failed to start")
		}
	}()
	<- ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	if err := worker.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed shutdown:%s", err)
	}

	wg.Wait()
}