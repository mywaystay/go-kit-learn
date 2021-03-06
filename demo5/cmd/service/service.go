package service

import (
	nethttp "net/http"
	"github.com/opentracing/opentracing-go"
	"github.com/go-kit/kit/log"
	"os"
	"net"
	"github.com/afex/hystrix-go/hystrix"
)


// http服务运行
func HTTPRun() {
	// log.NewNopLogger() 替换下面的log对象后，取消日志
	logger := log.NewLogfmtLogger(os.Stderr)
	tracer := opentracing.GlobalTracer()
	// 添加日志相关组件
	httpHandler := createHttpService(logger, tracer)
	logger.Log("demo5服务启动，http服务地址: ", "127.0.0.1:8092")
	// 开放9000端口给断路器使用
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	go func(handler *hystrix.StreamHandler) {
		var err chan error
		err <- nethttp.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	}(hystrixStreamHandler)

	err := nethttp.ListenAndServe(":8092", httpHandler)

	if nil != err {
		logger.Log("http监听发送错误：", err.Error())
	}
}

// grpc服务运行
func GRPCRun() {
	// log.NewNopLogger() 替换下面的log对象后，取消日志
	logger := log.NewLogfmtLogger(os.Stderr)
	tracer := opentracing.GlobalTracer()
	baseServer := createGRPCService(logger, tracer)
	// 监听9090端口
	grpcListener, err := net.Listen("tcp", ":9092")

	if nil != err {
		logger.Log("grpc监听发送错误：", err.Error())
	}

	logger.Log("demo5服务启动，grpc服务地址: ", "127.0.0.1:9092")
	err = baseServer.Serve(grpcListener)

	if nil != err {
		logger.Log("grpc服务启动错误：", err.Error())
	}

	grpcListener.Close()
}