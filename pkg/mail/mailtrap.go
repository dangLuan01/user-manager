package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/dangLuan01/user-manager/internal/utils"
)

type MailtrapConfig struct {
	MailSender string
	NameSender string
	MailTrapUrl string
	MailTrapApiKey string
}

type MailtrapProvider struct {
	client *http.Client
	config *MailtrapConfig
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {

	mailtrapConfig, ok := config.ProviderConfig["mailtrap"].(map[string]any)
	if !ok {
		return nil, utils.NewError(string(utils.ErrCodeInternal), "Invalid or missing MailTrap config")
	}


	return &MailtrapProvider{
		client: &http.Client{Timeout: config.Timeout},
		config: &MailtrapConfig{
			MailSender: mailtrapConfig["mail_sender"].(string),
			NameSender: mailtrapConfig["name_sender"].(string),
			MailTrapUrl: mailtrapConfig["mailtrap_url"].(string),
			MailTrapApiKey: mailtrapConfig["mailtrap_api_key"].(string),
		},
	}, nil
}

func (p *MailtrapProvider) SendMail(ctx context.Context, email *Email) error {

	email.From = Address{
		Email: p.config.MailSender,
		Name: p.config.NameSender,
	}

	payload, err := json.Marshal(email)
	if err != nil {
		log.Println(err)
		return utils.WrapError(string(utils.ErrCodeInternal), "Failed to marshal email", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.config.MailTrapUrl, bytes.NewReader(payload))
	if err != nil {
		log.Println(err)
		return utils.WrapError(string(utils.ErrCodeInternal), "Failed to create request", err)
	}

	req.Header.Add("Authorization", "Bearer " + p.config.MailTrapApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		log.Println(err)
		return utils.WrapError(string(utils.ErrCodeInternal), "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		return utils.NewError(string(utils.ErrCodeInternal), "Unexpected response from mailtrap")
	}

	return nil
}