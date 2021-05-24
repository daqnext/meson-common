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
	unauthorized    = 101
	failure         = 102
	malParams       = 103
	tokenError      = 104
	userForbidden   = 105
	lowerVersion    = 106
	captchaError    = 107
	captchaCoolDown = 108
	unknown         = 999

	//dns
	hostNotExist = 112

	//cdnuser request
	bindnameNotExist    = 201
	binddomainNotActive = 202
	notEnoughBalance    = 203
	fileNameError       = 204

	//register
	usernameExist    = 2001
	emailExist       = 2002
	phoneExist       = 2003
	vcodeError       = 2004
	usertypeError    = 2005
	emialFormatError = 2006

	//login
	usernameNotExist = 2101
	emailNotExist    = 2102
	phoneNotExist    = 2103
	emailError       = 2104
	phoneError       = 2105
	mismatchPwd      = 2004

	// /store/upload
	fileExist        = 5001
	noAliveFileStore = 5002

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

	addDownloadTaskFailed = 4001
	noEnoughSpace         = 4002
)

var (
	//common errorcode
	ErrUserUnAuth      = newHTTPErr(unauthorized, "user unauthorized")
	ErrInternalError   = newHTTPErr(failure, "server internal errorcode")
	ErrUnknown         = newHTTPErr(unknown, "unknown errorcode")
	ErrTokenError      = newHTTPErr(tokenError, "user token error")
	ErrUserForbidden   = newHTTPErr(userForbidden, "user forbidden")
	ErrLowerVersion    = newHTTPErr(lowerVersion, "your version need upgrade")
	ErrCaptcha         = newHTTPErr(captchaError, "captcha wrong")
	ErrCaptchaCoolDown = newHTTPErr(captchaCoolDown, "captcha cooldown")
	ErrMalParams       = newHTTPErr(malParams, "malformed request params")

	//DNS
	ErrHostNotExist = newHTTPErr(hostNotExist, "host not exist")

	//cdnuser request
	ErrBindNameNotExist    = newHTTPErr(bindnameNotExist, "bind name not exist")
	ErrBindDomainNotActive = newHTTPErr(binddomainNotActive, "bind domain not active")
	ErrNotEnoughBalance    = newHTTPErr(notEnoughBalance, "not enough balance")
	ErrFileNameError       = newHTTPErr(fileNameError, "file name error")

	//register
	ErrUsernameExist    = newHTTPErr(usernameExist, "username already exist")
	ErrEmailExist       = newHTTPErr(emailExist, "email already exist")
	ErrPhoneExist       = newHTTPErr(phoneExist, "phone already exist")
	ErrVcodeError       = newHTTPErr(vcodeError, "verification code error")
	ErrUserTypeError    = newHTTPErr(usertypeError, "usertype error")
	ErrEmailFormatError = newHTTPErr(emialFormatError, "email format error")

	//login
	ErrUsernameNotExist = newHTTPErr(usernameNotExist, "username not exist")
	ErrEmailNotExist    = newHTTPErr(emailNotExist, "email not exist")
	ErrPhoneNotExist    = newHTTPErr(phoneNotExist, "phone not exist")
	ErrPhoneError       = newHTTPErr(phoneError, "phone error")
	ErrEmailError       = newHTTPErr(emailError, "email error")
	ErrPwd              = newHTTPErr(mismatchPwd, "username or password is wrong")

	// /store/upload
	ErrFileExist = newHTTPErr(fileExist, "file already exist")
	ErrNoMachine = newHTTPErr(noAliveFileStore, "no alive machine")

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
	ErrNoSpace               = newHTTPErr(noEnoughSpace, "not enough space")
	ErrAddDownloadTaskFailed = newHTTPErr(addDownloadTaskFailed, "add download task failed")
)
