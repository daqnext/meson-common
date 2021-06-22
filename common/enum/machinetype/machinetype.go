package machinetype

type EMachine string

const (
	Terminal     EMachine = "Terminal"
	SpeedTester  EMachine = "SpeedTester"
	FileTransfer EMachine = "FileTransfer"
	Validator    EMachine = "Validator"
	FileStore    EMachine = "FileStore"
	RegionServer EMachine = "RegionServer"
	CenterServer EMachine = "CenterServer"
	LiveServer   EMachine = "LiveServer"
)
