package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/apiserver"
	"github.com/xyctruth/profiler/pkg/collector"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/storage/badger"
	"github.com/xyctruth/profiler/pkg/utils"
)

var (
	configPath     string
	dataPath       string
	dataGCInternal time.Duration
	uiGCInternal   time.Duration
)

func main() {
	flag.StringVar(&configPath, "config-path", "./collector.yaml", "Collector configuration file path")
	flag.StringVar(&dataPath, "data-path", "./data", "Collector Data file path")
	flag.DurationVar(&dataGCInternal, "data-gc-internal", 5*time.Minute, "Collector Data gc internal")
	flag.DurationVar(&uiGCInternal, "ui-gc-internal", 2*time.Minute, "Trace and pprof ui gc internal, must be greater than or equal to 1m")

	flag.Parse()

	log.WithFields(log.Fields{"configPath": configPath, "dataPath": dataPath, "dataGCInternal": dataGCInternal.String(), "uiGCInternal": uiGCInternal.String()}).
		Info("flag parse")

	if uiGCInternal < time.Minute {
		log.Fatal("ui-gc-internal must be greater than or equal to 1m")
		return
	}

	// Register the pprof endpoint
	utils.RegisterPProf()

	// New Store
	store := badger.NewStore(badger.DefaultOptions(dataPath).WithGCInternal(dataGCInternal))
	// Run collector
	collectorManger := runCollector(configPath, store)
	// Run api server
	apiServer := runAPIServer(store, uiGCInternal)

	// receive signal exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-quit
	log.Info("signal receive exit ", s)
	collectorManger.Stop()
	apiServer.Stop()
	store.Release()
}

// runAPIServer Run apis ,pprof ui ,trace ui
func runAPIServer(store storage.Store, gcInternal time.Duration) *apiserver.APIServer {
	apiServer := apiserver.NewAPIServer(
		apiserver.DefaultOptions(store).
			WithAddr(":8080").
			WithGCInternal(gcInternal))

	apiServer.Run()
	return apiServer
}

// runCollector Run collector manger
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
