package commands

type Commands struct {
	GetTipicalData      GetTipicalDataCmdHandler
	EventGet            EventGetCmdHandler
	PostTipicalMutation PostTipicalMutationCmdHandler
}

func NewCommands(
	getTipicalData GetTipicalDataCmdHandler,
	eventGet EventGetCmdHandler,
	postTipicalMutation PostTipicalMutationCmdHandler,
) *Commands {
	return &Commands{
		GetTipicalData:      getTipicalData,
		EventGet:            eventGet,
		PostTipicalMutation: postTipicalMutation,
	}
}
