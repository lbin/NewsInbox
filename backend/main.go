package main

import (
	"context"
	"fmt"
	"learnerai/conf"
	appContext "learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	mgodb "learnerai/models/mgodb/common"
	// milvus "learnerai/models/milvus/common"
	"learnerai/router"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func init() {
	ctx := appContext.GetGinContextWithRequestId()
	//viper初始化conf里边的配置文件
	conf.InitViper()
	//初始化日志信息
	initLogger()
	logger.Infof(ctx, "community init config execute start %d", time.Now().Unix())

	conf.GetGlobalConf(ctx)
	// mongodb
	mgodb.Setup()
	// milvus
	// milvus.Setup()
	initCrontab()
	// 设置redis,用于用户登录做缓存
	// client.SetupByHost(conf.GlobalEdgeSetting.Redis.Host, conf.GlobalEdgeSetting.Redis.Password)
	logger.Infof(ctx, "community init config execute end %d", time.Now().Unix())
}

func main() {
	ctx := appContext.GetGinContextWithRequestId()
	gin.SetMode(conf.ServerSetting.RunMode)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// go engine.StartKafka()

	mainCtx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	wg := wait(func() {
		defer cancel()
		if err := runHttpServer(mainCtx, ctx); err != nil {
			logger.Error(ctx, "mqtt server run got err", zap.Error(err))
		}
	})

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-mainCtx.Done():
		logger.Infof(ctx, "why run here!")
	case sg := <-quit:
		logger.Infof(ctx, "Receive signal %v and shutdown...", sg)
		//if inK8s {
		// k8s 中, 在接收到SIGINT 信号时, 仍然可能有一些延迟的流量打进来, 所以, 在k8s 中可以延迟关闭系统.
		logger.Infof(ctx, "delay cancel in %d seconds", conf.ServerSetting.ShutDownTimeout)
		time.Sleep(conf.ServerSetting.ShutDownTimeout)
		//}
		cancel()
	}
	wg.Wait()
	logger.Info(ctx, "All shutdown, exit")
}

func initCrontab() {
	crontab := cron.New()
	// _, _ = crontab.AddFunc("* * * * *", timer.ClearOldRecord) //清理录制的数据，每分钟一次
	crontab.Start()
}

func wait(funcs ...func()) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	wg.Add(len(funcs))
	for _, fn := range funcs {
		fn := fn
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	return wg
}

func runHttpServer(ctx context.Context, appContext *gin.Context) error {
	routersInit := router.InitRouter()
	readTimeout := conf.ServerSetting.ReadTimeout
	writeTimeout := conf.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", conf.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	svr := &http.Server{
		Addr:           endPoint,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	stop := make(chan error)
	go func() {
		logger.Infof(appContext, "Start http server listening %s", endPoint)
		if err := svr.ListenAndServe(); err == nil || err == http.ErrServerClosed {
			stop <- nil
		} else {
			stop <- errors.Wrap(err, "server serve err")
		}
	}()

	select {
	case <-ctx.Done():
		logger.Warnf(appContext, "Canceled, stop server %s", endPoint)
		cctx, cancel := context.WithTimeout(context.TODO(), conf.ServerSetting.ShutDownTimeout)
		defer cancel()
		return errors.Wrap(svr.Shutdown(cctx), "server shutdown err")
	case err := <-stop:
		return err
	}
}

func initLogger() {
	//获取日志的存储路径
	filePath := fmt.Sprintf("%s%s",
		conf.AppSetting.RuntimeRootPath,
		conf.AppSetting.LogSavePath,
	)

	fileName := fmt.Sprintf("%s.%s",
		conf.AppSetting.LogSaveName,
		conf.AppSetting.LogFileExt,
	)
	//获取配置的日志级别
	logLevel := conf.AppSetting.LogLevel
	runMode := conf.ServerSetting.RunMode
	expireDay := conf.AppSetting.ExpireTime
	logger.Setup(filePath, fileName, logLevel, runMode, int32(expireDay))
}
