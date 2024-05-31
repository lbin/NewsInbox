package request

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"learnerai/golib/common"
	uuid "github.com/satori/go.uuid"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code      int32       `json:"code"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"request_id"`
	Time      time.Time   `json:"time"`
	Data      interface{} `json:"data"`
}

func (g *Gin) Response(httpCode, errCode int32, data interface{}) {
	context.SetErrorCode(g.C, errCode)
	response := Response{
		Code:      errCode,
		Data:      data,
		Time:      time.Now(),
		RequestId: context.GetRequestID(g.C),
		Msg:       common.MsgFlags[errCode],
	}
	res, err := json.Marshal(response)
	if err != nil {
		logger.Errorf(g.C, "smart_request_out: marshal error [%s] ", err.Error())
	}
	logger.Debugf(g.C, "smart_request_out: res[%s] ", string(res))
	g.C.JSON(int(httpCode), response)
	return
}

func (g *Gin) ResponseWithMessage(httpCode, errCode int32, data interface{}, msg string) {
	context.SetErrorCode(g.C, errCode)
	response := Response{
		Code:      errCode,
		Msg:       msg,
		Data:      data,
		Time:      time.Now(),
		RequestId: context.GetRequestID(g.C),
	}
	res, err := json.Marshal(response)
	if err != nil {
		logger.Errorf(g.C, "smart_request_out: marshal error [%s] ", err.Error())
	}
	logger.Debugf(g.C, "smart_request_out: res[%s] ", string(res))
	g.C.JSON(int(httpCode), response)
	return
}

func (g *Gin) ResponseWithAny(httpCode, errCode int32, data interface{}) {
	context.SetErrorCode(g.C, errCode)

	response := map[string]interface{}{
		"time":       time.Now(),
		"request_id": context.GetRequestID(g.C),
	}

	if data != nil {
		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		if t.Kind() != reflect.Struct {
			logger.Errorf(g.C, "ResponseWithAny data type error, must be struct")
			return
		}

		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")

			switch v.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				response[tag] = v.Field(i).Int()
			case reflect.Float32, reflect.Float64:
				response[tag] = v.Field(i).Float()
			case reflect.String:
				response[tag] = v.Field(i).String()
			}
		}
	}

	logger.Debugf(g.C, "smart_request_out: res[%+v] ", data)

	g.C.JSON(int(httpCode), response)
}

func (g *Gin) ResponseSuccess(data interface{}) {
	g.Response(http.StatusOK, common.SUCCESS, data)
}

func (g *Gin) ResponseError(errCode int32, data interface{}) {
	g.Response(http.StatusOK, errCode, data)
}

func (g *Gin) ResponseErrorWithMsg(errCode int32, msg string) {
	g.ResponseWithMessage(http.StatusOK, errCode, nil, msg)
}

func (g *Gin) ResponseErrorWithMsgAndData(errCode int32, data interface{}, msg string) {
	g.ResponseWithMessage(http.StatusOK, errCode, data, msg)
}

func (g *Gin) BindAndValid(form interface{}) error {
	err := g.C.ShouldBind(form)
	if err != nil {
		logger.Warn(g.C, err.Error())
		return err
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		logger.Warn(g.C, err.Error())
		return err
	}
	if !check {
		logger.Warn(g.C, "valid check failed")
		return valid.Errors[0]
	}

	return nil
}

func GetRequestId() (requestId string) {
	value := uuid.Must(uuid.NewV4(), nil).String()
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func AddRequestId(ctx *gin.Context) {
	requestId := GetRequestId()
	context.SetRequestID(ctx, requestId)
	body, err := ctx.GetRawData()
	if err != nil {
		logger.Error(ctx, err.Error())
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//输入参数写入日志
	logger.Infof(ctx, "smart_request_in: RequestURI[%s],RequestBody[%v]", ctx.Request.Header, string(body))
}

func AddRequestIdWithoutLog(ctx *gin.Context) {
	requestId := GetRequestId()
	context.SetRequestID(ctx, requestId)
	body, err := ctx.GetRawData()
	if err != nil {
		logger.Error(ctx, err.Error())
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func (g *Gin) ResponseWithMessageZh(errCode int32, err error) {
	logger.Errorf(g.C, "errCode[%+v],err[%+v]", errCode, err)

	errMsg := ""
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			errMsg = err.Translate(*Trans)
			break
		}
	}

	if errMsg == "" {
		errMsg = common.GetMsg(errCode)
	}

	g.ResponseWithMessage(http.StatusOK, errCode, nil, errMsg)
}
