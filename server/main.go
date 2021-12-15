package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/felixge/fgprof"
	log "github.com/sirupsen/logrus"
	"github.com/xyctruth/profiler/pkg/apiserver"
	"github.com/xyctruth/profiler/pkg/collector"
	"github.com/xyctruth/profiler/pkg/storage"
	"github.com/xyctruth/profiler/pkg/storage/badger"
)

var (
	configPath  string
	storagePath string
)

func main() {
	if configPath = os.Getenv("config_path"); configPath == "" {
		configPath = "./collector.yaml"
	}

	if storagePath = os.Getenv("storage_path"); storagePath == "" {
		storagePath = "./data"
	}

	log.Info("configPath:", configPath, " storagePath:", storagePath)

	registerPProf()

	store := badger.NewStore(storagePath)
	collectorManger := runCollector(configPath, store)
	apiServer := runAPIServer(store)

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
	go apiServer.Run()
	return apiServer
}

func runCollector(configPath string, store storage.Store) *collector.Manger {
	m := collector.NewManger(store)
	collector.LoadConfig(configPath, func(config collector.Config) {
		log.Info("config change, reload collector!!!")
		m.Load(config)
	})

	return m
}

func registerPProf() {
	go func() {
		http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			panic(err)
		}
	}()

}
