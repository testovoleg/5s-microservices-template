package commands

type Commands struct {
	GetTipicalData GetTipicalDataCmdHandler

	PostTipicalMutation PostTipicalMutationCmdHandler
}

func NewCommands(
	getTipicalData GetTipicalDataCmdHandler,

	postTipicalMutation PostTipicalMutationCmdHandler,
) *Commands {
	return &Commands{
		GetTipicalData: getTipicalData,

		PostTipicalMutation: postTipicalMutation,
	}
}
