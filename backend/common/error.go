package common

var ErrNoDocuments = "mongo: no documents in result"

// GetMsg get error information based on Code
func GetMsg(code int32) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

const (
	SUCCESS                     = 200   //成功
	ERROR                       = 500   //失败
	InvalidParams               = 400   //参数错误
	SSONotLoggedIn              = 1100  // sso未登录
	NotLoggedIn                 = 1000  // 未登录
	ParameterIllegal            = 1001  // 参数不合法
	UnauthorizedUserId          = 1002  // 用户Id不合法
	Unauthorized                = 1003  // 未授权
	ServerError                 = 1004  // 系统错误
	NotData                     = 1005  // 没有数据
	ModelAddError               = 1006  // 添加错误
	ModelDeleteError            = 1007  // 删除错误
	ModelStoreError             = 1008  // 存储错误
	OperationFailure            = 1009  // 操作失败
	RoutingNotExist             = 1010  // 路由不存在
	ErrorUserExist              = 1011  // 用户已存在
	ErrorUserNotExist           = 1012  // 用户不存在
	ErrorNeedCaptcha            = 1013  // 需要验证码
	PasswordInvalid             = 1014  // 密码不符合规范
	SignError                   = 1015  // 签名错误
	ErrorCheckTokenFail         = 10001 // 用户信息获取失败
	ErrorCheckUserRoleFail      = 10002 // 用户信息获取失败
	ErrorUploadSaveImageFail    = 50000 //图片上传失败
	ErrorUploadCheckImageFormat = 50001 //格式错误
	ErrorUploadCheckImageFail   = 50002 //格式错误
	PhoneInvalid                = 10007
	LicenseInvalid              = 10008
	CustomIdentityInvalid       = 10010 // 身份信息不正确
	AccessTooFrequently         = 99999 // 访问太频繁
)

var MsgFlags = map[int32]string{
	SUCCESS:                  "成功",
	ERROR:                    "失败",
	InvalidParams:            "参数不合法",
	ErrorCheckTokenFail:      "用户校验失败，请尝试重新登录（10001）",
	ErrorCheckUserRoleFail:   "用户不存在",
	ErrorUploadSaveImageFail: "上传图片失败",
	NotLoggedIn:              "没有登录",
	ParameterIllegal:         "参数错误",
	UnauthorizedUserId:       "用户验证失败",
	Unauthorized:             "用户验证失败",
	ErrorNeedCaptcha:         "用户验证失败",
	NotData:                  "无数据",
	ServerError:              "服务器错误",
	ModelAddError:            "添加失败",
	ModelDeleteError:         "删除失败",
	ModelStoreError:          "存储失败",
	OperationFailure:         "操作失败",
	RoutingNotExist:          "路由不存在",
	ErrorUserExist:           "用户已存在",
	ErrorUserNotExist:        "用户不存在",
	PasswordInvalid:          "密码不符合规范",
	SignError:                "签名错误",
	SSONotLoggedIn:           "用户未登录",
	CustomIdentityInvalid:    "身份信息不正确",
	AccessTooFrequently:      "访问太频繁",
}
