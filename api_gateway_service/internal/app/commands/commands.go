package commands

type Commands struct {
	ApiMethods ApiMethodsCmdHandler

	WebhookMethods WebhookMethodsCmdHandler
}

func NewCommands(
	apiMethods ApiMethodsCmdHandler,

	webhookMethods WebhookMethodsCmdHandler,
) *Commands {
	return &Commands{
		ApiMethods: apiMethods,

		WebhookMethods: webhookMethods,
	}
}
