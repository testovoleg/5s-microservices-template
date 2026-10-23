package commands

type Commands struct {
	ApiMethods     ApiMethodsCmdHandler
	CommonMethods  CommonMethodsCmdHandler
	WebhookMethods WebhookMethodsCmdHandler
}

func NewCommands(
	apiMethods ApiMethodsCmdHandler,
	commonMethods CommonMethodsCmdHandler,
	webhookMethods WebhookMethodsCmdHandler,
) *Commands {
	return &Commands{
		ApiMethods:     apiMethods,
		CommonMethods:  commonMethods,
		WebhookMethods: webhookMethods,
	}
}
