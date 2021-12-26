package main

import (
	"flag"
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
	flag.StringVar(&configPath, "config-path", "./collector.yaml", "Collector configuration file path")
	flag.StringVar(&storagePath, "data-path", "./data", "The path to store the collected data")
	flag.Parse()

	log.Info("configPath:", configPath, " storagePath:", storagePath)

	// Register the pprof endpoint
	utils.RegisterPProf()

	store := badger.NewStore(storagePath)
	collectorManger := runCollector(configPath, store)
	apiServer := runAPIServer(store)

	// receive signal exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-quit
	log.Info("signal receive exit ", s)
	collectorManger.Stop()
	apiServer.Stop()
	store.Release()
}

func runAPIServer(store storage.Store) *apiserver.APIServer {
	apiServer := apiserver.NewAPIServer(":8081", store)
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
