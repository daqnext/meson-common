package commonmsg

type MachineStateBaseMsg struct {
	MacAddr       string  `json:"mac_addr"`
	MemTotal      uint64  `json:"mem_total"` // uint: byte
	MemAvailable  uint64  `json:"mem_avail"`
	DiskTotal     uint64  `json:"disk_total"`
	DiskAvailable uint64  `json:"disk_avail"`
	Version       string  `json:"version"`
	CpuUsage      float64 `json:"cpu_usage"`
}

type TerminalStatesMsg struct {
	OS               string `json:"os"`
	CPU              string `json:"cpu"` // cpu model name
	Port             string `json:"port"`
	CdnDiskTotal     uint64 `json:"cdn_disk_total"`
	CdnDiskAvailable uint64 `json:"cdn_disk_avail"`
	MachineStateBaseMsg
}

type FileTransferStateMsg struct {
	MachineStateBaseMsg
}

type SpeedTesterStateMsg struct {
	MachineStateBaseMsg
}

type ValidatorStateMsg struct {
	MachineStateBaseMsg
}

type FileStoreStateMsg struct {
	MachineStateBaseMsg
}

type DownLoadFileCmdMsg struct {
	DownloadUrl  string `json:"downloadurl" binding:"required"`
	TransferTag  string `json:"transfertag" binding:"required"`
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	FileSize     int64  `json:"filesize"  binding:"required"`
	Continent    string `json:"continent" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Area         string `json:"area" binding:"required"`
}

type IpfsUploadUrlMsg struct {
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	DownloadUrl  string `json:"download_url" binding:"required"`
}

type FileTransferDownLoadFinishMsg struct {
	FilePartHash string `json:"fileparthash" binding:"required"`
	FileSize     int64  `json:"filesize" binding:"required"`
	TransferTag  string `json:"transfertag" binding:"required"`
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	Continent    string `json:"continent" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Area         string `json:"area" binding:"required"`
}

type FileTransferDownLoadFailedMsg struct {
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	Continent    string `json:"continent" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Area         string `json:"area" binding:"required"`
}

type TerminalDownloadFinishMsg struct {
	TransferTag  string `json:"transfertag" binding:"required"`
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	Continent    string `json:"continent" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Area         string `json:"area" binding:"required"`
}

type TerminalDownloadFailedMsg struct {
	TransferTag  string `json:"transfertag" binding:"required"`
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	Continent    string `json:"continent" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Area         string `json:"area" binding:"required"`
}

type TerminalRequestDeleteFilesMsg struct {
	Files []string `json:"files" binding:"required"`
}

type FileTransferDeleteFileMsg struct {
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
}

type SpeedReportMsg struct {
	MachineTag string `json:"machine_tag" binding:"required"`
	Speed      uint64 `json:"speed"`
}

type SpeedTestCmdMsg struct {
	MachineTag     string `json:"machine_tag" binding:"required"`
	Port           string `json:"port" binding:"required"`
	FileName       string `json:"file_name" binding:"required"`
	DownloadSecond int    `json:"download_second" binding:"required"`
}

type DeleteFolderCmdMsg struct {
	FolderName string `json:"foldername" binding:"required"`
}

type UserUploadToFileStoreFinish struct {
	UploadUserIp string `json:"uploadUserIp" binding:"required"`
	//OriginUrl string `json:"originUrl" binding:"required"`
	Size           int64  `json:"size" binding:"required"`
	UserName       string `json:"userName" binding:"required"`
	OriginFileName string `json:"originFileName"`
	FileName       string `json:"fileName" binding:"required"`
	FileHash       string `json:"fileHash" binding:"required"`
	FileSystem     string `json:"fileSystem" binding:"required"`
}

type DBRecordMsg struct {
	ARecordMap     map[string][4]byte
	TxtRecordArray [][]string
	CaaRecordArray []string
	NSRecordArray  []string
	CNameRecordMap map[string]string
}

type ValidateStruct struct {
	TerminalTag string
	BindName    string
	FileName    string
	PartHash    string
}

type ValidateFailMsg struct {
}

type RedisConnectionDataMsg struct {
	Host           string
	Port           int
	Auth           string
	MaxPoolSize    int
	MaxIdle        int
	IdleTimeoutSec int
	Db             int
}
