package gb

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

var (
	ShortMsgClient = new(ShortMsg)
)

type ShortMsg struct {
	client *dysmsapi20170525.Client
}

func InitShortMsgClient(credentialConfig *credential.Config, endpoint string) error {
	newCredential, err := credential.NewCredential(credentialConfig)
	if err != nil {
		return err
	}

	config := &openapi.Config{
		Credential: newCredential,
	}
	config.Endpoint = tea.String(endpoint)
	result := &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(config)
	if err != nil {
		return err
	}
	ShortMsgClient = &ShortMsg{
		client: result,
	}
	return nil
}

func InitShortMsgSimpleClient(accessKeyId, accessKeySecret string) error {
	return InitShortMsgClient(new(credential.Config).
		SetType("access_key").
		SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret),
		"dysmsapi.aliyuncs.com")
}

func (s *ShortMsg) SendMsg(sendSmsRequest *dysmsapi20170525.SendSmsRequest) error {
	runtime := &util.RuntimeOptions{}
	tryErr := func() (e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				e = r
			}
		}()
		_, err := s.client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		return tryErr
	}
	return nil
}

func (s *ShortMsg) SendSimpleMsg(targetPhoneNumber, signName, templateCode, templateParam string) error {
	return s.SendMsg(&dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(targetPhoneNumber),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	})
}
