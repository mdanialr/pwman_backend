package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mdanialr/pwman_backend/internal/app"
	"github.com/mdanialr/pwman_backend/internal/middleware"
	conf "github.com/mdanialr/pwman_backend/pkg/config"
	gormLogger "github.com/mdanialr/pwman_backend/pkg/gorm"
	help "github.com/mdanialr/pwman_backend/pkg/helper"
	"github.com/mdanialr/pwman_backend/pkg/postgresql"
	"github.com/mdanialr/pwman_backend/pkg/storage"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HTTP setup and serve Http server using fiber app.
func HTTP() {
	// init viper config
	v, err := conf.InitConfigYml()
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}

	// init zap logger
	zapLog, err := setupZapLogger(v)
	if err != nil {
		log.Fatalln("failed to init zap app logger:", err)
	}
	defer zapLog.Sync()

	// init GORM log writer
	gormLogWr, err := setupGORMWriter(v)
	if err != nil {
		log.Fatalln("failed to init gorm logger:", err)
	}
	defer gormLogWr.Close()
	// init GORM logger
	gormLog := gormLogger.New(gormLogWr, gormLogger.WithLogLevel(v.GetInt("gorm.lvl")))
	// init GORM with some customization
	db, err := postgresql.NewGorm(v,
		postgresql.WithCustomLogger(gormLog),
		postgresql.WithDisableNestedTrx(),
		postgresql.WithPrepareStatement(),
		postgresql.WithSkipDefaultTrx(),
		postgresql.WithSingularTableName(),
	)
	if err != nil {
		log.Fatalln("failed to init gorm:", err)
	}

	// init storage provider
	st := storage.NewFile(zapLog)

	// init fiber app log writer
	fiberLogWr, err := setupFiberWriter(v)
	if err != nil {
		log.Fatalln("failed to init fiber app logger:", err)
	}
	defer fiberLogWr.Close()
	// init fiber log middleware
	fbLog := fiberLog.New(fiberLog.Config{
		Output:     fiberLogWr,
		TimeFormat: time.DateTime,
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/metrics" {
				return true
			}
			return false
		},
	})
	// set default value for metrics refresh rate
	monRefreshRate := v.GetInt64("metrics.refresh")
	if monRefreshRate == 0 {
		monRefreshRate = 2 // set default to 2 seconds
	}
	// init fiber monitor metrics config
	monConf := setupFiberMetricsMonitor(v)
	// conditionally add proxy header from Nginx
	var proxyHeader string
	if v.GetString("server.env") == "prod" {
		proxyHeader = "X-Real-Ip"
	}
	// init fiber app
	fiberApp := fiber.New(fiber.Config{
		ProxyHeader:           proxyHeader,
		ReadTimeout:           10 * time.Second,
		IdleTimeout:           5 * time.Second,
		BodyLimit:             v.GetInt("server.limit") * 1024 * 1024,
		RequestMethods:        []string{fiber.MethodHead, fiber.MethodGet, fiber.MethodPost},
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		ErrorHandler:          help.DefaultHTTPErrorHandler,
		DisableStartupMessage: v.GetString("server.env") == "prod", // disable startup message on production env
	})
	// add useful middlewares for fiber app
	fiberApp.Use(
		fbLog,
		recover.New(),
		compress.New(),
	)
	// assign metrics to endpoint
	fiberApp.Get("/metrics",
		// add guard to this endpoint, since this endpoint will expose hardware resource info
		func(c *fiber.Ctx) error {
			if c.Query("pass") == v.GetString("metrics.pass") {
				return c.Next()
			}
			return c.SendStatus(fiber.StatusNotFound)
		},
		monitor.New(monConf),
	)
	// server file in /dl
	dir := strings.TrimSuffix(v.GetString("storage.path"), "/")
	fiberApp.Use("/dl",
		// give jwt middleware before accessing any media resources
		middleware.JWT(v),
		filesystem.New(filesystem.Config{Root: http.Dir(dir)}),
	)

	// init internal http handlers
	h := app.HttpHandler{
		R:       fiberApp.Group("/api"),
		DB:      db,
		Config:  v,
		Log:     zapLog,
		Storage: st,
	}
	h.SetupRouter()

	// log the app host and port
	host := v.GetString("server.host") + ":" + v.GetString("server.port")
	zapLog.Info("Run app in " + host)

	// listen from a different goroutine
	go func() {
		if err := fiberApp.Listen(host); err != nil {
			log.Panicf("failed listen into port %v", err)
		}
	}()

	// create channel for signal being sent
	c := make(chan os.Signal, 1)
	// when an interrupt or termination signal is sent, notify the channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// blocks main thread until an interrupt is received
	<-c
	zapLog.Info("gracefully shutting down...")
	fiberApp.Shutdown()
	zapLog.Info("running cleanup tasks...")
	// some clean up task should be done here
	zapLog.Sync()
	zapLog.Info("services was successful shutdown.")
}

func setupZapLogger(conf *viper.Viper) (*zap.Logger, error) {
	var zapConfig zap.Config

	// determine which zap env should be used
	env := conf.GetString("server.env")
	switch env {
	case "prod":
		zapConfig = zap.NewProductionConfig()
	case "dev":
		zapConfig = zap.NewDevelopmentConfig()
	default:
		return nil, errors.New("unsupported env. only support prod and dev")
	}

	// determine which output should be used by zap
	logType := conf.GetString("zap.log")
	switch logType {
	case "console":
		zapConfig.Encoding = "console"
	case "json":
		zapConfig.Encoding = "json"
		logPath := strings.TrimSuffix(conf.GetString("zap.path"), "/") + "/log"
		// make sure the output log path is not empty
		if logPath == "" {
			return nil, errors.New("zap.path is required when zap.log is json")
		}
		zapConfig.OutputPaths = []string{logPath}
	default:
		return nil, errors.New("unsupported zap encoding type. only support console and json")
	}

	return zapConfig.Build()
}

func setupGORMWriter(conf *viper.Viper) (*os.File, error) {
	// determine which output should be used by GORM logger
	logType := conf.GetString("gorm.log")
	switch logType {
	case "console":
		return os.Stdout, nil
	case "file":
		logPath := conf.GetString("gorm.path")
		// make sure the output log path is not empty
		if logPath == "" {
			return nil, errors.New("gorm.path is required when gorm.log is file")
		}
		// set target file log
		targetLog := strings.TrimSuffix(logPath, "/") + "/gorm-log"
		return os.OpenFile(targetLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	default:
		return nil, errors.New("unsupported gorm logger output. only support console and file")
	}
}

func setupFiberWriter(conf *viper.Viper) (*os.File, error) {
	// determine which output should be used by GORM logger
	logType := conf.GetString("fiber.log")
	switch logType {
	case "console":
		return os.Stdout, nil
	case "file":
		logPath := conf.GetString("fiber.path")
		// make sure the output log path is not empty
		if logPath == "" {
			return nil, errors.New("fiber.path is required when fiber.log is file")
		}
		// set target file log
		targetLog := strings.TrimSuffix(logPath, "/") + "/fiber-app-log"
		return os.OpenFile(targetLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	default:
		return nil, errors.New("unsupported gorm logger output. only support console and file")
	}
}

func setupFiberMetricsMonitor(conf *viper.Viper) monitor.Config {
	// set default value for metrics refresh rate
	refRate := conf.GetInt64("metrics.refresh")
	if refRate == 0 {
		refRate = 2 // set default to 2 seconds
	}
	// set default monitor title
	title := conf.GetString("metrics.title")
	if title == "" {
		title = "Password Manager API Monitor"
	}
	return monitor.Config{
		Title:   title,
		Refresh: time.Duration(refRate) * time.Second,
	}
}
