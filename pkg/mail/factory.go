package mail

import (
	"fmt"

	"github.com/dangLuan01/user-manager/internal/utils"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type ProviderFactory interface {
	CreateProvider(config *MailConfig) (EmailProviderService, error)
}

type MailtrapProviderFactory struct {

}

func (f *MailtrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error)  {
	return NewMailtrapProvider(config)
}

func NewProviderFactory(providerType ProviderType) (ProviderFactory, error) {
	switch providerType {
	case ProviderMailtrap:
		return &MailtrapProviderFactory{}, nil
	default:
		return nil, utils.NewError(string(utils.ErrCodeInternal),fmt.Sprintf("Unsupported provider type:%s", providerType))
	}
	
}