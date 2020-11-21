package resp

type code int

type httpErr struct {
	code code
	msg  string
}

func (e *httpErr) Code() code {
	return e.code
}

func (e *httpErr) Error() string {
	return e.msg
}

func newHTTPErr(code code, msg string) *httpErr {
	return &httpErr{
		code: code,
		msg:  msg,
	}
}

const (
	success code = 0
	//common errorcode
	unauthorized  = 101
	failure       = 102
	malParams     = 103
	tokenError    = 104
	userForbidden = 105
	unknown       = 999

	//cdnuser request
	bindnameNotExist    = 201
	binddomainNotActive = 202
	notEnoughBalance    = 203

	//register
	usernameExist = 2001
	emailExist    = 2002
	phoneExist    = 2003
	vcodeError    = 2004
	usertypeError = 2005

	//login
	usernameNotExist = 2101
	emailNotExist    = 2102
	phoneNotExist    = 2103
	mismatchPwd      = 2004

	// client/newdomain
	bindnameExist  = 1001
	originurlExist = 1002

	// client/deletetdomain
	domainDeleteNotExist = 1201

	// client/modifydomain
	domainModifyNotExist = 1101

	// t/bindname/*action
	fileLinkExpired = 2201

	// ================= terminal part =================
	saveFileFailed = 3001
	setIndexFailed = 3002
	fileNotExist   = 3003

	// ================= FileTransfer part =================
	addDownloadTaskFailed = 4001
)

var (
	//common errorcode
	ErrUserUnAuth    = newHTTPErr(unauthorized, "user unauthorized")
	ErrInternalError = newHTTPErr(failure, "server internal errorcode")
	ErrUnknown       = newHTTPErr(unknown, "unknown errorcode")
	ErrTokenError    = newHTTPErr(tokenError, "user token error")
	ErrUserForbidden = newHTTPErr(userForbidden, "user forbidden")
	ErrMalParams     = newHTTPErr(malParams, "malformed request params")

	//cdnuser request
	ErrBindNameNotExist    = newHTTPErr(bindnameNotExist, "bind name not exist")
	ErrBindDomainNotActive = newHTTPErr(binddomainNotActive, "bind domain not active")
	ErrNotEnoughBalance    = newHTTPErr(notEnoughBalance, "not enough balance")

	//register
	ErrUsernameExist = newHTTPErr(usernameExist, "username already exist")
	ErrEmailExist    = newHTTPErr(emailExist, "email already exist")
	ErrPhoneExist    = newHTTPErr(phoneExist, "phone already exist")
	ErrVcodeError    = newHTTPErr(vcodeError, "verification code error")
	ErrUserTypeError = newHTTPErr(usertypeError, "usertype error")

	//login
	ErrUsernameNotExist = newHTTPErr(usernameNotExist, "username not exist")
	ErrEmailNotExist    = newHTTPErr(emailNotExist, "email not exist")
	ErrPhoneNotExist    = newHTTPErr(phoneNotExist, "phone not exist")
	ErrPwd              = newHTTPErr(mismatchPwd, "username or password is wrong")

	// client/newdomain
	ErrBindnameExist  = newHTTPErr(bindnameExist, "bindname already exist")
	ErrOriginurlExist = newHTTPErr(originurlExist, "originurl already exist")

	// client/deletedomain
	ErrDomainDeleteNotExist = newHTTPErr(domainDeleteNotExist, "domain id not exist")

	// client/modifydomain
	ErrDomainModifyNotExist = newHTTPErr(domainModifyNotExist, "domain id not exist")

	// t/bindname/*action
	ErrFileLinkExpired = newHTTPErr(fileLinkExpired, "file link expired")

	// ================= terminal part =================
	ErrSaveFile     = newHTTPErr(saveFileFailed, "failed to save file")
	ErrSetFileIndex = newHTTPErr(setIndexFailed, "failed to set index for new file")
	ErrFileNotExist = newHTTPErr(fileNotExist, "file not exist in local index")

	// ================= FileTransfer part =================
	ErrAddDownloadTaskFailed = newHTTPErr(addDownloadTaskFailed, "add download task failed")
)
