package machinetype

type EMachine string

const (
	Terminal     EMachine = "Terminal"
	SpeedTester  EMachine = "SpeedTester"
	FileTransfer EMachine = "FileTransfer"
	Validator    EMachine = "Validator"
)
