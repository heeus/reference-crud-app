/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

const ServiceVersion = "0.0.2"

const (
	mnPutCount        = "putCount"
	mnBatchCount      = "batchCount"
	mnFlush           = "batchDuration"
	mnPartitionsSize  = "partitionsSize"
	mnHcCnt           = "hcCnt"
	mnHcDurNs         = "hcDurNs"
	mnBatchInterval   = "batchInterval"
	mnCacheViewCnt    = "cacheViewCnt"
	mnNotCacheViewCnt = "notCacheViewCnt"
)

//Service s.e.
type Service struct {
	driver DBDriver
	port   int64

	readFunc   string
	insertFunc string
	updateFunc string
	deleteFunc string
	scanFunc   string

	pathPattern string

	timestart int64

	noop bool

	logger *Logger

	EventCount      int64
	BatchCount      int64
	BatchDurationNS int64
	PartitionsSize  int64
	HcCnt           int64
	HcDurNs         int64
	CacheViewCnt    int64
	NotCacheViewCnt int64
}

//Init s.e.
func (s *Service) Init() error {
	args := mapArgs(os.Args)

	s.noop = initBoolParam(args, NoopServiceEnvironmentProperty, NoopServiceAttribute, false)

	s.logger = &Logger{}
	s.logger.level = initIntParam(args, LoggerLevelEnvironmentProperty, LoggerLevelAttribute, 0)

	s.flushMetrics()

	if s.noop {
		s.driver = &NopDriver{logger: s.logger}
	} else if d, err := s.getServiceDriver(args); err == nil {
		s.driver = d
	} else {
		s.logger.Error(err.Error())
		return err
	}

	if p, err := s.getServicePort(args); err == nil {
		s.port = p
	} else {
		s.logger.Error(err.Error())
		return err
	}

	if err := s.driver.Init(args); err != nil {
		s.logger.Error(err.Error())
		return err
	}

	//init path pattern

	s.pathPattern = initStringParam(args, PathPatternEnvironmentProperty, PathPatternAttribute, DefaultPathPattern)
	s.readFunc = initStringParam(args, ServiceReadFuncEnvironmentProperty, ServiceReadFuncAttribute, ReadDefaultFunc)
	s.insertFunc = initStringParam(args, ServiceInsertFuncEnvironmentProperty, ServiceInsertFuncAttribute, InsertDefaultFunc)
	s.updateFunc = initStringParam(args, ServiceUpdateFuncEnvironmentProperty, ServiceUpdateFuncAttribute, UpdateDefaultFunc)
	s.deleteFunc = initStringParam(args, ServiceDeleteFuncEnvironmentProperty, ServiceDeleteFuncAttribute, DeleteDefaultFunc)
	s.scanFunc = initStringParam(args, ServiceScanFuncEnvironmentProperty, ServiceScanFuncAttribute, ScanDefaultFunc)

	s.logger.Debug("service successfully initialized")

	return nil
}

//Start s.e.
func (s *Service) Start() {
	defer s.Stop()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(s.Handle404)

	r.HandleFunc(s.pathPattern, s.handle)
	r.HandleFunc(s.pathPattern+"/", s.handle)

	r.HandleFunc("/api/driver/clean", s.handleClean)
	r.HandleFunc("/api/driver/clean/", s.handleClean)

	r.HandleFunc("/api/vars", s.handleVars)
	r.HandleFunc("/api/vars/", s.handleVars)

	r.HandleFunc("/api", s.handleRoot)
	r.HandleFunc("/api/", s.handleRoot)

	//todo add port arg

	path := fmt.Sprintf(":%v", s.port)

	go func() {
		s.logger.Log("Listening at localhost%v \n", path)

		s.timestart = time.Now().Unix()

		http.ListenAndServe(path, r)
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}

//Stop  s.e.
func (s *Service) Stop() {
	s.flushMetrics()
	s.driver.Free()
	s.logger.Log("Service stoped")
}

func (s *Service) handleMetrics(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{}

	resp["Status"] = 200
	resp["Error"] = ""

	resp[mnBatchCount] = s.getBatchCount()
	resp[mnFlush] = s.getBatchDuration()
	resp[mnBatchInterval] = 0
	resp[mnHcCnt] = s.getMetricHcCnt()
	resp[mnHcDurNs] = s.getMetricHcDurNs()
	resp[mnPartitionsSize] = s.getPartitionsSize()
	resp[mnPutCount] = s.getPutCount()
	resp[mnCacheViewCnt] = s.getCacheViewCnt()
	resp[mnNotCacheViewCnt] = s.getNotCacheViewCnt()

	bytes, err := json.Marshal(resp)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)

	//s.flushMetrics()

}

func (s *Service) handleClean(w http.ResponseWriter, r *http.Request) {
	res := s.driver.Clean(nil)

	if res.Error != "" {
		s.logger.Error("DB driver clean error: %v", res.Error)
	}

	bytes, err := json.Marshal(res)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

//Handle404 s.e.
func (s *Service) Handle404(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("Service asked for not supported route: %q", r.URL.Path)
	fmt.Fprintf(w, "Service can't route this path (404 not fount); \n\n")
}

func (s *Service) handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service works: \n\n")
	fmt.Fprintf(w, "Service version: %v\n", ServiceVersion)
	fmt.Fprintf(w, "Driver: %v. \n", s.driver.Name())
	fmt.Fprintf(w, "Port: %v. \n", s.port)
	fmt.Fprintf(w, "Starts at: %v. \n", time.Unix(s.timestart, 0))
	fmt.Fprintf(w, "Path pattern: %v. \n", s.pathPattern)
	fmt.Fprintf(w, "\n %v", s.driver.Info())
}

func (s *Service) handleVars(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(r, []string{"POST", "GET"}); err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintf(w, "Enviroment variables: \n\n")
	envs := os.Environ()

	if len(envs) > 0 {

		for k, e := range envs {
			fmt.Fprintf(w, "%v. %v;\n", k, e)
		}
	}
}

func (s *Service) handle(w http.ResponseWriter, r *http.Request) {
	startHc := time.Now()

	if s.noop {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("{\"Status\":200}"))
		return
	}

	params := mux.Vars(r)

	f := params["function"]

	if f == "YcsbMetric" {
		s.handleMetrics(w, r)
		return
	}

	wsid := params["wsid"]

	req, err := buildRequest(r)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Debug("Handle request: ")
	s.logger.Debug("%v", r.URL.Path)
	s.logger.Debug("%v", req)

	req.Partition, err = strconv.ParseInt(wsid, 10, 64)

	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res *DBResponse

	startBatch := time.Now()

	switch f {
	case s.readFunc:
		s.logger.Debug("Processing read request")
		res = s.driver.Read(req)
	case s.insertFunc:
		res = s.driver.Insert(req)
	case s.updateFunc:
		res = s.driver.Update(req)
	case s.scanFunc:
		res = s.driver.Scan(req)
	case s.deleteFunc:
		res = s.driver.Delete(req)
	default:
		str := fmt.Sprintf("Func %q not allowed!", f)
		s.logger.Error(str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	atomic.AddInt64(&s.BatchDurationNS, time.Since(startBatch).Nanoseconds())

	if res.Error != "" {
		s.logger.Error("DB driver proccessing error: %v", res.Error)
	}

	bytes := res.stringify()

	s.logger.Debug("Response: %v", string(bytes))

	atomic.AddInt64(&s.HcDurNs, time.Since(startHc).Nanoseconds())

	atomic.AddInt64(&s.EventCount, 1)
	atomic.AddInt64(&s.HcCnt, 1)
	atomic.AddInt64(&s.BatchCount, 1)
	atomic.AddInt64(&s.CacheViewCnt, 0)
	atomic.AddInt64(&s.NotCacheViewCnt, 1)

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func (s *Service) getServicePort(args map[string]string) (int64, error) {
	port := initIntParam(args, ServicePortEnvironmentProperty, ServicePortAttribute, DefaultPort)

	if port < 0 && port > 65535 {
		return 0, fmt.Errorf("port value should be greater than 0 and less than 65535")
	}

	s.logger.Debug("Selected service port: %v", port)

	return port, nil
}

func (s *Service) getServiceDriver(args map[string]string) (driver DBDriver, err error) {
	driverName := initStringParam(args, ServiceDriverEnvironmentProperty, ServiceDriverAttribute, "cas")

	s.logger.Debug("Selected db driver: %q", driverName)

	switch driverName {
	case "cas":
		return &CasandraDriver{logger: s.logger}, nil
	case "casp":
		return &CasandraPartitionedDriver{logger: s.logger}, nil
	case "light":
		return &LightDriver{logger: s.logger}, nil
	case "mem":
		return &MemoryDriver{logger: s.logger}, nil
	default:
		return nil, fmt.Errorf("wrong driver is given. Available: cas, light, mem")
	}
}

func (s *Service) flushMetrics() {
	atomic.StoreInt64(&s.EventCount, 0)
	atomic.StoreInt64(&s.BatchCount, 0)
	atomic.StoreInt64(&s.BatchDurationNS, 0)
	atomic.StoreInt64(&s.PartitionsSize, 1)
	atomic.StoreInt64(&s.HcCnt, 0)
	atomic.StoreInt64(&s.HcDurNs, 0)
	atomic.StoreInt64(&s.CacheViewCnt, 0)
	atomic.StoreInt64(&s.NotCacheViewCnt, 0)
}

func (s *Service) getPutCount() int64 {
	return atomic.LoadInt64(&s.EventCount)
}

func (s *Service) getCacheViewCnt() int64 {
	return atomic.LoadInt64(&s.CacheViewCnt)
}

func (s *Service) getNotCacheViewCnt() int64 {
	return atomic.LoadInt64(&s.NotCacheViewCnt)
}

func (s *Service) getBatchCount() int64 {
	return atomic.LoadInt64(&s.BatchCount)
}

func (s *Service) getBatchDuration() int64 {
	return atomic.LoadInt64(&s.BatchDurationNS)
}

func (s *Service) getPartitionsSize() int64 {
	return atomic.LoadInt64(&s.PartitionsSize)
}

func (s *Service) getMetricHcCnt() int64 {
	return atomic.LoadInt64(&s.HcCnt)
}

func (s *Service) getMetricHcDurNs() int64 {
	return atomic.LoadInt64(&s.HcDurNs)
}
