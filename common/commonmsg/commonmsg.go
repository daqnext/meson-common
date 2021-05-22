package commonmsg

import "github.com/daqnext/meson-common/common/enum/machinetype"

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
	OS               string     `json:"os"`
	CPU              string     `json:"cpu"` // cpu model name
	Port             string     `json:"port"`
	NetInMbs         [5]float64 `json:"net_in_mbs"`
	NetOutMbs        [5]float64 `json:"net_out_mbs"`
	CdnDiskTotal     uint64     `json:"cdn_disk_total"`
	CdnDiskAvailable uint64     `json:"cdn_disk_avail"`
	MachineSetupTime string     `json:"machine_setup_time"`
	SequenceId       int        `json:"sequence_id"`
	MachineStateBaseMsg
}

type TerminalInfoMsg struct {
	OS               string `json:"os"`
	CPU              string `json:"cpu"` // cpu model name
	Port             string `json:"port"`
	MachineSetupTime string `json:"machine_setup_time"`
	MacAddr          string `json:"mac_addr"`
}

type TerminalHeatBeatMsg struct {
	MemTotal         uint64  `json:"mem_total"` // uint: byte
	MemAvailable     uint64  `json:"mem_avail"`
	DiskTotal        uint64  `json:"disk_total"`
	DiskAvailable    uint64  `json:"disk_avail"`
	CdnDiskTotal     uint64  `json:"cdn_disk_total"`
	CdnDiskAvailable uint64  `json:"cdn_disk_avail"`
	NetInMbs         float64 `json:"net_in_mbs"`
	NetOutMbs        float64 `json:"net_out_mbs"`
	Version          string  `json:"version"`
	CpuUsage         float64 `json:"cpu_usage"`
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

type SignMsg struct {
	TimeStamp  int64  `json:"timestamp"`
	MachineMac string `json:"mac"`
	Sign       string `json:"sign"`
	MacSign    string `json:"mac_sign"`
}

type TransferPauseMsg struct {
	PauseTime int `json:"pausetime"`
	SignMsg
}

type CrossRegionTransferFileCmdMsg struct {
	DownloadUrl      string `json:"downloadurl" binding:"required"`
	BindName         string `json:"bindname" binding:"required"`
	FileName         string `json:"filename" binding:"required"`
	FileSize         uint64 `json:"filesize"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	OriginRegion     string `json:"originRegion"`
}

type CrossRegionTransferFinishMsg struct {
	BindName             string `json:"bindname" binding:"required"`
	FileName             string `json:"filename" binding:"required"`
	FileSize             uint64 `json:"filesize"`
	RequestContinent     string `json:"continent"`
	RequestCountry       string `json:"country"`
	RequestArea          string `json:"area"`
	OriginRegion         string `json:"origin_region"`
	CachedRegion         string `json:"cached_region"`
	TransferTerminalPort string `json:"transfer_terminal_port"`
	TransferTerminalTag  string `json:"transfer_terminal_tag"`
}

type CrossRegionTransferFailedMsg struct {
	BindName         string `json:"bindname" binding:"required"`
	FileName         string `json:"filename" binding:"required"`
	FileSize         uint64 `json:"filesize"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	OriginRegion     string `json:"origin_region"`
	FailedRegion     string `json:"failed_region"`
}

type IpfsUploadUrlMsg struct {
	BindNameHash string `json:"bindnamehash" binding:"required"`
	FileNameHash string `json:"filenamehash" binding:"required"`
	DownloadUrl  string `json:"download_url" binding:"required"`
}

type DownLoadFileCmdMsg struct {
	DownloadUrl      string `json:"downloadurl" binding:"required"`
	BindName         string `json:"bindname" binding:"required"`
	FileName         string `json:"filename" binding:"required"`
	FileSize         uint64 `json:"filesize"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	DownloadType     string `json:"downloadType"`
	OriginRegion     string `json:"originRegion"`
	//TargetRegion     string `json:"targetRegion"`
	SignMsg
}

type DeleteFileCmdMsg struct {
	BindName string `json:"bindname" binding:"required"`
	FileName string `json:"filename" binding:"required"`
	SignMsg
}

type TerminalDownloadFinishMsg struct {
	BindName         string `json:"bindname" binding:"required"`
	FileName         string `json:"filename" binding:"required"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	DownloadType     string `json:"downloadType"`
	OriginRegion     string `json:"originRegion"`
	//TargetRegion     string `json:"targetRegion"`
	DownloadUrl string `json:"downloadUrl"`
	FileSize    uint64 `json:"filesize"`
}

type TerminalDownloadFailedMsg struct {
	BindName         string `json:"bindname" binding:"required"`
	FileName         string `json:"filename" binding:"required"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	DownloadType     string `json:"downloadType"`
	OriginRegion     string `json:"originRegion"`
	//TargetRegion     string `json:"targetRegion"`
	DownloadUrl string `json:"downloadUrl"`
	FileSize    uint64 `json:"filesize"`
}

type TerminalDownloadProcessMsg struct {
	BindName         string `json:"bindname"`
	FileName         string `json:"filename"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
	Downloaded       int64  `json:"downloaded"`
}

type TerminalDownloadStartMsg struct {
	BindName         string `json:"bindname"`
	FileName         string `json:"filename"`
	RequestContinent string `json:"continent"`
	RequestCountry   string `json:"country"`
	RequestArea      string `json:"area"`
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
	MacAddr        string `json:"machine_mac" binding:"required"`
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

type PanicReportMsg struct {
	MachineType machinetype.EMachine `json:"machineType" binding:"required"`
	TimeStamp   int64
	Region      string
	TerminalStatesMsg
	Error string `json:"error" binding:"required"`
	Stack string `json:"stack" binding:"required"`
}
