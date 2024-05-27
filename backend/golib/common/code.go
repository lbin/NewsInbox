package common

import "github.com/pkg/errors"

const (
	SUCCESS                       = 200   //成功
	ERROR                         = 500   //失败
	InvalidParams                 = 400   //参数错误
	NotLoggedIn                   = 1000  // 未登录
	SSONotLoggedIn                = 1100  // sso未登录
	ParameterIllegal              = 1001  // 参数不合法
	UnauthorizedUserId            = 1002  // 用户Id不合法
	Unauthorized                  = 1003  // 未授权
	ServerError                   = 1004  // 系统错误
	NotData                       = 1005  // 没有数据
	ModelAddError                 = 1006  // 添加错误
	ModelDeleteError              = 1007  // 删除错误
	ModelStoreError               = 1008  // 存储错误
	OperationFailure              = 1009  // 操作失败
	RoutingNotExist               = 1010  // 路由不存在
	ErrorUserExist                = 1011  // 用户已存在
	ErrorUserNotExist             = 1012  // 用户不存在
	ErrorNeedCaptcha              = 1013  // 需要验证码
	PasswordInvalid               = 1014  //密码不符合规范
	PrivateHostError              = 1016  // 边缘服务器私有域名不存在
	FPMSAddOrderError             = 1017  // 添加工单失败
	FPMSFinishOrderError          = 1018  // 关闭工单失败
	FPMSGroupTreeError            = 1019  // 获取项目树失败
	MaxQueryTimeError             = 1020  // 查询时间达到上限
	NoPermissionsForCommunity     = 1021  //没有小区查看权限
	CheckFaceQualityFailed        = 1022  //人脸质量检测失败
	TimeConfigOverlaps            = 1023  // 时间重叠
	ErrorCheckTokenFail           = 10001 //用户信息获取失败
	ErrorCheckUserRoleFail        = 10002 // 用户信息获取失败
	ErrorReportDeviceInsertFail   = 20001 //上报设备信息失败
	ErrorDeviceUpdateFail         = 20002 // 设备更新失败
	ErrorReportDeviceUpsertFail   = 20003 //上报设备信息失败
	ErrorGetTrafficMapFail        = 20004 //获取人流地图数据失败
	ErrorDeviceDeleteFail         = 20005 //设备删除失败
	ErrorReportDeviceMonitorFail  = 20006 // 算法上报数据失败
	ErrorAlarmInsertFail          = 30001 //报警上报失败
	ErrorAlarmDealFail            = 30002 //报警处理失败
	ErrorAddPostConfigFail        = 40001 //查岗配置添加失败
	ErrorAddPatrolConfigFail      = 40002 //巡更点配置失败
	ErrorUpdatePatrolConfigFail   = 40003 //巡更更新失败
	ErrorAddPatrolRecordFail      = 40004 //巡查失败
	ErrorNoPatrolConfig           = 40005 //没有巡查点配置
	ErrorSendWarningCheckAuthFail = 40007 // 一键报警验证失败
	ErrorUploadSaveImageFail      = 50000 //图片上传失败
	ErrorUploadCheckImageFormat   = 50001 //格式错误
	ErrorUploadCheckImageFail     = 50002 //格式错误
	ErrorRTSPUrlInvalid           = 10003 //视频流地址错误
	LicenseNotExist               = 10004 //车牌号不存在
	ReachMaxPlayCount             = 10005 //到达播放上限
	LicenseExist                  = 10006 //车牌号不存在
	VeichlePropertyInvalid        = 10009 // 车辆属性不正确
	CustomIdentityInvalid         = 10010 // 身份信息不正确
	AccessTooFrequently           = 99999 // 访问太频繁
	PhoneInvalid                  = 10007
	LicenseInvalid                = 10008
	TaskNameExist                 = 10011
)

var MsgFlags = map[int32]string{
	SUCCESS:                       "成功",
	ERROR:                         "失败",
	InvalidParams:                 "参数不合法",
	ErrorCheckTokenFail:           "用户校验失败，请尝试重新登录（10001）",
	ErrorCheckUserRoleFail:        "用户不存在",
	ErrorReportDeviceInsertFail:   "上报设备信息失败",
	ErrorReportDeviceUpsertFail:   "上报设备信息失败",
	ErrorGetTrafficMapFail:        "获取人流地图数据失败",
	ErrorDeviceUpdateFail:         "设备更新失败",
	ErrorDeviceDeleteFail:         "设备更新失败",
	ErrorAlarmInsertFail:          "报警上报失败",
	ErrorAlarmDealFail:            "报警上报失败",
	ErrorAddPostConfigFail:        "添加查岗配置失败",
	ErrorAddPatrolConfigFail:      "添加巡更配置失败",
	ErrorUpdatePatrolConfigFail:   "更新巡更配置失败",
	ErrorAddPatrolRecordFail:      "巡更失败",
	ErrorUploadSaveImageFail:      "上传图片失败",
	ErrorReportDeviceMonitorFail:  "摄像头数据上报失败",
	ErrorSendWarningCheckAuthFail: "一键报警验证失败",
	NotLoggedIn:                   "没有登录",
	ParameterIllegal:              "参数错误",
	UnauthorizedUserId:            "用户验证失败",
	Unauthorized:                  "用户验证失败",
	ErrorNeedCaptcha:              "用户验证失败",
	NotData:                       "无数据",
	ServerError:                   "服务器错误",
	ModelAddError:                 "添加失败",
	ModelDeleteError:              "删除失败",
	ModelStoreError:               "存储失败",
	OperationFailure:              "操作失败",
	RoutingNotExist:               "路由不存在",
	ErrorUserExist:                "用户已存在",
	ErrorUserNotExist:             "用户不存在",
	ErrorRTSPUrlInvalid:           "rtsp流地址错误",
	LicenseNotExist:               "车牌号不存在",
	ReachMaxPlayCount:             "播放视频到达上限",
	LicenseExist:                  "车牌号已存在",
	PasswordInvalid:               "密码不符合规范",
	VeichlePropertyInvalid:        "车辆属性不正确",
	CustomIdentityInvalid:         "身份信息不正确",
	AccessTooFrequently:           "访问太频繁",
	ErrorNoPatrolConfig:           "没有巡查点配置",
	ErrorUploadCheckImageFormat:   "上传文件格式错误",
	ErrorUploadCheckImageFail:     "上传文件失败",
	PhoneInvalid:                  "手机号格式错误",
	LicenseInvalid:                "车牌号格式错误",
	PrivateHostError:              "边缘服务器私有域名不存在",
	SSONotLoggedIn:                "用户未登录",
	FPMSAddOrderError:             "添加工单失败",
	FPMSFinishOrderError:          "关闭工单失败",
	FPMSGroupTreeError:            "获取项目树失败",
	MaxQueryTimeError:             "最多只能查询一个月的数据哦",
	NoPermissionsForCommunity:     "您没有查看此小区的权限或小区不存在",
	CheckFaceQualityFailed:        "您上传的人脸图片不合格，请按要求重新上传！",
	TimeConfigOverlaps:            "任务时段重复，请修改",
	TaskNameExist:                 "此任务名称重复，请更换任务名称",
}

var MongoNoFound = errors.New("MONGO NO FOUND")
var MongoInsertError = errors.New("MONGO INSERT ERROR")
var DeviceNoFound = errors.New("DEVICE NO FOUND")
var FileNoFound = errors.New("FILE NO FOUND")
var StatusError = errors.New("STATUS ERROR")
var NoNeedInsert = errors.New("NO NEED INSERT")
var NoData = errors.New("NO DATA")

// GetMsg get error information based on Code
func GetMsg(code int32) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
