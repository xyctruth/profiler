package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/apiserver"
	"github.com/xyctruth/profiler/pkg/collector"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/storage/badger"
	"github.com/xyctruth/profiler/pkg/utils"
)

var (
	configPath  string
	storagePath string
)

func main() {
	if configPath = os.Getenv("CONFIG_PATH"); configPath == "" {
		configPath = "./collector.yaml"
	}

	if storagePath = os.Getenv("DATA_PATH"); storagePath == "" {
		storagePath = "./data"
	}

	log.Info("configPath:", configPath, " storagePath:", storagePath)

	utils.RegisterPProf()
	store := badger.NewStore(storagePath)
	collectorManger := runCollector(configPath, store)
	apiServer := runAPIServer(store)

	// receive signal exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	log.Info("signal receive exit ", s)
	collectorManger.Stop()
	apiServer.Stop()
	store.Release()
}

func runAPIServer(store storage.Store) *apiserver.APIServer {
	apiServer := apiserver.NewAPIServer(":8080", store)
	apiServer.Run()
	return apiServer
}

func runCollector(configPath string, store storage.Store) *collector.Manger {
	m := collector.NewManger(store)
	err := collector.LoadConfig(configPath, func(config collector.CollectorConfig) {
		log.Info("config change, reload collector!!!")
		m.Load(config)
	})
	if err != nil {
		panic(err)
	}
	return m
}
