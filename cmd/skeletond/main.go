package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"

	"github.com/pocoz/skeleton/config"
	"github.com/pocoz/skeleton/db/elasticsearch"
	"github.com/pocoz/skeleton/db/mssql"
	"github.com/pocoz/skeleton/db/postgres"
	"github.com/pocoz/skeleton/db/postgres/migrator"
	"github.com/pocoz/skeleton/models"
	"github.com/pocoz/skeleton/services/httpserver"
	"github.com/pocoz/skeleton/services/memlogger"
	"github.com/pocoz/skeleton/services/scheduler"
	"github.com/pocoz/skeleton/services/scrollsvc"
)

func main() {
	const (
		exitCodeSuccess = 0
		exitCodeFailure = 1
	)

	errc := make(chan error, 1)
	donec := make(chan struct{})
	sigc := make(chan os.Signal, 1)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "app", config.ServiceName, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfg, err := config.New()
	if err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(exitCodeFailure)
	}

	svcMemLogger := memlogger.New(&memlogger.Config{
		Logger: logger,
	})
	go svcMemLogger.Run()

	engineES, err := elasticsearch.New(&elasticsearch.Config{
		Server:   cfg.ElasticServer,
		User:     cfg.ElasticUser,
		Password: cfg.ElasticPassword,
		Context:  ctx,
		Logger:   logger,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize elasticsearch", "err", err)
		os.Exit(exitCodeFailure)
	}

	engineMSSQL, err := mssql.New(&mssql.Config{
		Server:   cfg.SQLServer,
		Port:     cfg.SQLPort,
		User:     cfg.SQLUser,
		Password: cfg.SQLPassword,
		Database: cfg.SQLDatabase,
		Context:  ctx,
		Logger:   logger,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize mssql", "err", err)
		os.Exit(exitCodeFailure)
	}

	credentialsPostgres := &models.CredentialsDB{
		Server:   cfg.PostgresServer,
		Port:     cfg.PostgresPort,
		User:     cfg.PostgresUser,
		Password: cfg.PostgresPassword,
		Database: cfg.PostgresDatabase,
	}

	migratorPostgres := migrator.New(&migrator.Config{
		Context:       ctx,
		Logger:        logger,
		CredentialsDB: credentialsPostgres,
		Folder:        cfg.PostgresMigrationsFolder,
		Table:         cfg.PostgresMigrationsTable,
	})
	err = migratorPostgres.Run()
	if err != nil {
		level.Error(logger).Log("msg", "migrations was failure", "err", err)
		os.Exit(exitCodeFailure)
	}

	storePostgres, err := postgres.New(&postgres.Config{
		Context:       ctx,
		Logger:        logger,
		MaxButchSize:  cfg.PostgresMaxButchSize,
		CredentialsDB: credentialsPostgres,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize postgres", "err", err)
		os.Exit(exitCodeFailure)
	}

	scrollSvc, err := scrollsvc.New(&scrollsvc.Config{
		Logger:  logger,
		DBES:    engineES,
		Context: ctx,
		ConfigTasks: &scrollsvc.ConfigTasks{
			TaskTemplateName:        cfg.TaskTemplateName,
			TaskTemplateCronPattern: cfg.TaskTemplateCronPattern,
			TaskTemplateIsActive:    cfg.TaskTemplateIsActive,
		},
	})

	schedulerSVC, err := scheduler.New(logger)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize scheduler svc", "err", err)
		os.Exit(exitCodeFailure)
	}

	err = schedulerSVC.Run(scrollSvc.Tasks)
	if err != nil {
		level.Error(logger).Log("msg", "failed to run scheduler", "err", err)
		os.Exit(exitCodeFailure)
	}

	httpServer, err := httpserver.New(&httpserver.Config{
		Port:        cfg.HTTPPort,
		Logger:      logger,
		RateLimiter: rate.NewLimiter(rate.Every(cfg.RateLimitEvery), cfg.RateLimitBurst),
		DBES:        engineES,
		DBMsSQL:     engineMSSQL,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize http server", "err", err)
		os.Exit(exitCodeFailure)
	}
	go func() {
		level.Info(logger).Log("msg", "starting http server", "port", cfg.HTTPPort)
		err = httpServer.Run()
		if err != nil {
			level.Error(logger).Log("msg", "http server run failure", "err", err)
			os.Exit(exitCodeFailure)
		}
	}()

	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	go func() {
		select {
		case sig := <-sigc:
			level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)
			httpServer.Shutdown()   // Stop http server
			schedulerSVC.Shutdown() // Stop all running tasks
			engineMSSQL.Shutdown()  // Close mssql connection
			engineES.Shutdown()     // Close elasticsearch connection
			storePostgres.Shutdown()
			signal.Stop(sigc)
			close(donec)
		case <-errc:
			level.Info(logger).Log("msg", "now exiting with error", "error code", exitCodeFailure)
			os.Exit(exitCodeFailure)
		}
	}()

	<-donec
	level.Info(logger).Log("msg", "goodbye")
	os.Exit(exitCodeSuccess)
}
