package commands

type Commands struct {
	AddApi     AddApiCmdHandler
	GetApi     GetApiCmdHandler
	GetFullApi GetFullApiCmdHandler
	UpdateApi  UpdateApiCmdHandler
	DeleteApi  DeleteApiCmdHandler

	Webhook WebhookCmdHandler
}

func NewCommands(
	addApi AddApiCmdHandler,
	getApi GetApiCmdHandler,
	getFullApi GetFullApiCmdHandler,
	updateApi UpdateApiCmdHandler,
	deleteApi DeleteApiCmdHandler,

	webhook WebhookCmdHandler,
) *Commands {
	return &Commands{
		AddApi:     addApi,
		GetApi:     getApi,
		GetFullApi: getFullApi,
		UpdateApi:  updateApi,
		DeleteApi:  deleteApi,

		Webhook: webhook,
	}
}
