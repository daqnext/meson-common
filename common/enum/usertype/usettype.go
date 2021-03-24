package usertype

type EUser string

const (
	Client       EUser = "client"
	Terminal     EUser = "terminal"
	Admin        EUser = "admin"
	Validator    EUser = "validator"
	FileTransfer EUser = "filetransfer"
	SpeedTester  EUser = "speedtester"
	Blog         EUser = "blog"
	FileStore    EUser = "filestore"
	CenterServer EUser = "centerserver"
	RegionServer EUser = "regionserver"
)
