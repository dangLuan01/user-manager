package mail

import (
	"context"
	"log"
	"time"

	"github.com/dangLuan01/user-manager/internal/config"
	"github.com/dangLuan01/user-manager/internal/utils"
)

type Email struct {
	From 		Address 	`json:"from"`
	To 			[]Address 	`json:"to"`
	Subject 	string 		`json:"subject"`
	Text 		string 		`json:"text"`
	Category 	string 		`json:"category"`
}

type Address struct {
	Email 	string `json:"email"`
	Name 	string `json:"name,omitempty"`
}

type MailConfig struct {
	ProviderConfig 	map[string]any
	ProviderType 	ProviderType
	MaxRetries 		int
	Timeout 		time.Duration
}

type MailService struct {
	config *MailConfig
	provider EmailProviderService
}

func NewMailService(cfg *config.Config, providerFactory ProviderFactory) (EmailProviderService, error) {

	config := &MailConfig{
		ProviderConfig: cfg.MailProviderConfig,
		ProviderType: ProviderType(cfg.MailProviderType),
		MaxRetries: 3,
		Timeout: 10 * time.Second,
	}

	provider, err := providerFactory.CreateProvider(config)
	if err != nil {
		log.Println(err)
		return nil, utils.WrapError(string(utils.ErrCodeInternal), "Failed to create provider", err)
	}

	return &MailService{
		config: config,
		provider: provider,
	}, nil
}

func (ms *MailService) SendMail(ctx context.Context, email *Email) error {
	var lastErr error

	for attemps := 1; attemps < ms.config.MaxRetries; attemps++ {

		err := ms.provider.SendMail(ctx, email)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Println(lastErr)
	}

	log.Println(lastErr)
	return utils.WrapError(string(utils.ErrCodeInternal), "Failed to send email after retries", lastErr)
}