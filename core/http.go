package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	np "net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/middleware/requestid"

	"irir-layout/pkg/erroron"
	"irir-layout/pkg/log"
	"irir-layout/pkg/response"

	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/x/errors"

	"irir-layout/config"
	_ "irir-layout/docs"
)

// Router 加载路由，使用侧提供接口，实现侧需要实现该接口
type Router interface {
	Load(engine *iris.Application)
}

// HttpServer 通用 web 服务.
type HttpServer struct {
	BindAddress       string         // 监听地址
	BindPort          int            // 监听端口
	Debug             bool           // 启动模式
	CertKey           config.CertKey // 启用https
	Middlewares       []string       // 启用的中间件
	Health            bool           // 是否启用健康检查
	EnableMetrics     bool           // 是否启用监控
	EnableProfiling   bool           // 是否启用性能分析工具
	*iris.Application                // web 驱动
}

// NewHttpServer 从给定的配置返回 GenericAPIServer 的新实例。
func NewHttpServer(cnf *config.Config) *HttpServer {
	s := &HttpServer{
		BindAddress: cnf.HttpServerConfig.BindAddress,
		BindPort:    cnf.HttpServerConfig.BindPort,
		CertKey: config.CertKey{
			CertFile: cnf.HttpServerConfig.ServerCert.CertFile,
			KeyFile:  cnf.HttpServerConfig.ServerCert.KeyFile,
		},
		Debug:           cnf.HttpServerConfig.Debug,
		Health:          cnf.HttpServerConfig.Health,
		Middlewares:     cnf.HttpServerConfig.Middlewares,
		EnableMetrics:   cnf.HttpServerConfig.EnableMetrics,
		EnableProfiling: cnf.HttpServerConfig.EnableProfiling,
		Application:     iris.Default(),
	}

	InitGenericAPIServer(s)
	return s
}

// InitGenericAPIServer 初始化通用 API 服务
func InitGenericAPIServer(s *HttpServer) {
	if s.Debug {
		// 启动 API 文档
		s.SetupSwagger()
	}
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// address 将主机 IP 地址和主机端口号连接成一个地址字符串，例如：0.0.0.0:8443。
func (h *HttpServer) address() string {
	return net.JoinHostPort(h.BindAddress, strconv.Itoa(h.BindPort))
}

func (h *HttpServer) Setup() {
	if h.Debug {
		h.SetupSwagger()
	}

	h.Validator = validator.New()

	h.Application.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusInternalServerError)
		response.Error(ctx, erroron.ErrInternalServer, nil)
	})
}

// LoadRouter 加载自定义路由
func (h *HttpServer) LoadRouter(rs ...Router) {
	for _, r := range rs {
		r.Load(h.Application)
	}
}

func (h *HttpServer) InstallAPIs() {
	// 添加健康检查api
	if h.Health {
		h.Get("/health", func(ctx iris.Context) {
			_ = ctx.JSON(iris.Map{"code": 0, "msg": "OK"})
		})
	}

	// 添加监控
	// if h.EnableMetrics {
	// 	prometheus := ginprometheus.NewPrometheus("gin")
	// 	prometheus.Use(h.Engine)
	// }

	// 添加性能测试工具
	if h.EnableProfiling {
		h.Application.HandleMany("GET", "/debug/pprof /debug/pprof/{action:path}", pprof.New())
		h.Application.Get("/debug/pprof/profile", func(ctx iris.Context) {
			np.Profile(ctx.ResponseWriter(), ctx.Request())
		})
		h.Application.Get("/debug/pprof/symbol", func(ctx iris.Context) {
			np.Symbol(ctx.ResponseWriter(), ctx.Request())
		})
		h.Application.Get("/debug/pprof/cmdline", func(ctx iris.Context) {
			np.Cmdline(ctx.ResponseWriter(), ctx.Request())
		})
		h.Application.Get("/debug/pprof/trace", func(ctx iris.Context) {
			np.Trace(ctx.ResponseWriter(), ctx.Request())
		})
	}
}

// InstallMiddlewares 初始化中间件。
func (h *HttpServer) InstallMiddlewares() {
	// 必要中间件
	logConf := logger.DefaultConfig()
	logConf.LogFuncCtx = log.FuncCtx
	h.UseRouter(recover.New())
	h.UseRouter(requestid.New())
	h.UseRouter(logger.New(logConf))
	h.UseRouter(cors.New().Handler())
	h.UseRouter(iris.Compression)
}

// Run 启动 http 服务器.
func (h *HttpServer) Run() {
	quit := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 优雅退出
	iris.RegisterOnInterrupt(func() {
		quitErr := h.Application.Shutdown(ctx)
		if quitErr != nil {
			h.Logger().Errorf("quit error: %s", quitErr)
		}
		close(quit)
	})

	// 健康检查
	go func() {
		if h.Health {
			if err := h.ping(ctx); err != nil {
				h.Logger().Fatal(err.Error())
			}
		}
	}()

	h.Logger().Infof("Start to listening the incoming requests on http address: %s", h.address())
	var (
		key, cert = h.CertKey.KeyFile, h.CertKey.CertFile
		serverErr error
	)

	cf := config.GetConfig()

	if cert == "" || key == "" {
		serverErr = h.Listen(h.address(), iris.WithOptimizations, iris.WithConfiguration(cf.Iris))
	} else {
		serverErr = h.Application.Run(iris.TLS(h.address(), cert, key), iris.WithOptimizations, iris.WithConfiguration(cf.Iris))
	}

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		h.Logger().Fatal(serverErr.Error())
	}

	h.Logger().Infof("Server on %s stopped", h.address())
	<-quit
}

// ping 服务器健康
func (h *HttpServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/health", h.address())
	if strings.Contains(h.address(), "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%d/health", h.BindPort)
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			h.Logger().Info("The router has been deployed successfully.")
			_ = resp.Body.Close()

			return nil
		}

		// 暂停 1 秒钟
		h.Logger().Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			h.Logger().Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}

// SetupSwagger 启用swagger
//
//go:generate swag init -g ../cmd/main.go -o ../docs
func (h *HttpServer) SetupSwagger() {
	swaggerUI := swagger.Handler(swaggerFiles.Handler,
		swagger.URL("/swagger/swagger.json"),
		swagger.DeepLinking(true),
		swagger.Prefix("/swagger"),
	)
	h.Get("/swagger", swaggerUI)
	h.Get("/swagger/{any:path}", swaggerUI)
}
