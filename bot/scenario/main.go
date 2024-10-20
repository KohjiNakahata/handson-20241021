package main

import (
	"encoding/json"
	"flag"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/util"
	"handson/bot/scenario/lib/log"
	"handson/bot/scenario/lib/report"
	"handson/bot/scenario/scenarios"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var gp = &scenarios.GlobalParams{Host: "localhost:7000", ReceiveByteSize: 1400, UDPSendInterval: 100, LogLevel: 20}
var config string = "./bot/scenario/config"

type ScenarioSettings struct {
	ScenarioName    string `json:"type"`
	ScenarioPattern string `json:"run"`
	HowMany         int    `json:"howmany"`
	Duration        int    `json:"duration"`
	Interval        int    `json:"interval"`
}

var ss *ScenarioSettings
var scenarioFactory *func() scenarios.Scenario
var logger = log.New("BOT/MAIN")

func load() error {
	flag.StringVar(&ss.ScenarioName, "type", "Connect", "Scenario name that you implement and defined in ScenarioList.")
	flag.StringVar(&ss.ScenarioPattern, "run", "ConnectUDP", "Scenario instance that is defined as 'hint' in Json file.")
	flag.IntVar(&ss.HowMany, "howmany", -1, "The number of clients to join matching.")
	flag.IntVar(&ss.Duration, "duration", -1, "Duration to run the scenario in seconds.")
	flag.IntVar(&ss.Interval, "interval", -1, "Interval to create scenario clients in millisecond.")
	flag.Parse()
	return nil
}
func setup() error {
	gp.UserState = scenarios.NewUserState()
	scenarioFactory_, ok := scenarios.ScenarioFactoryList[ss.ScenarioName]
	if !ok {
		return util.NewError("No Scenario named \"%v\". Please check 'ScenarioList' in bot/scenario/scenarios/main.go", ss.ScenarioName)
	}
	scenarioFactory = &scenarioFactory_
	filepath.Walk(config, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Cannot read file. path:%v", path)
			return err
		}
		filename := info.Name()
		logger.Debug("Reading file \"%s\" ...", filename)
		if strings.HasSuffix(filename, ".json") {
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			logger.Debug("File Contents: %v", string(raw))
			json.Unmarshal(raw, &gp.Raw.CommonParams)
			scenarioParamsIf, ok := gp.Raw.CommonParams[ss.ScenarioPattern]
			if ok {
				gp.Raw.ScenarioParams = scenarioParamsIf.(map[string]any)
			}
		}
		return nil
	})
	if gp.Raw.ScenarioParams == nil {
		logger.Warn("No Scenario Parameter [%s] found in json files. Using only global parameters.", ss.ScenarioPattern)
	}
	globalParamsBytes, _ := gp.GenerateParams(0)
	err := json.Unmarshal(globalParamsBytes, &gp)
	logger.Sys("Parsed global params. %#v", gp)
	if err != nil {
		return util.StackError(util.NewError("Failed to pars common params. %v", err.Error()))
	}
	udp.LogLevel(gp.LogLevel)
	if gp.MetricsInterval > 0 {
		report.Interval = gp.MetricsInterval
	}
	if ss.HowMany >= 0 {
		gp.HowMany = ss.HowMany
	}
	if ss.Duration >= 0 {
		gp.Duration = ss.Duration
	}
	if ss.Interval >= 0 {
		gp.Interval = ss.Interval
	}
	logger.Info("Setup done. CommonParams:%+v ScenarioParams:%+v ParamsFromAPI:%+v", gp.Raw.CommonParams, gp.Raw.ScenarioParams, gp.Raw.ParamsFromAPI)
	return nil
}
func start() error {
	report.ResetAllMetrics()
	clients := make([]*scenarios.Scenario, gp.HowMany)
	var wg sync.WaitGroup
	wg.Add(gp.HowMany)
	for i := 0; i <= gp.HowMany-1; i++ {
		go func(i int) {
			defer wg.Done()
			scenarioClient := (*scenarioFactory)()
			clients[i] = &scenarioClient
			scenarioParamsBytes, _ := gp.GenerateParams(i)
			scenarioClient.ParseParam(i, scenarioParamsBytes)
			userID := scenarioClient.GetUserID()
			gp.UserState.Init(userID)
			err := scenarioClient.Run(gp)
			if err != nil {
				logger.Error(util.StackError(util.NewError("Scenario execution failed. %v", err.Error())))
				report.IncrementScenarioError()
				return
			}
			checkIdling := func() {
				startedAt := time.Now()
				lastActiveTime := time.Now()
				for startedAt.Add(time.Duration(gp.Duration) * time.Second).After(time.Now()) {
					if report.IsActive(userID) {
						lastActiveTime = time.Now()
					}
					if lastActiveTime.Add(time.Duration(gp.IdleDuration) * time.Second).Before(time.Now()) {
						logger.Verboseu(userID, "Triggering OnIdle... ")
						scenarioClient.OnIdle()
						lastActiveTime = time.Now()
					}
					time.Sleep(time.Second)
				}
			}
			go checkIdling()
		}(i)
		time.Sleep(time.Duration(gp.Interval) * time.Millisecond)
	}
	if gp.Duration == 0 {
		wg.Wait()
	} else {
		time.Sleep(time.Duration(gp.Duration) * time.Second)
	}
	for i := 0; i <= gp.HowMany-1; i++ {
		(*clients[i]).OnScenarioEnd()
	}
	report.StopAllMetrics()
	if gp.Duration < report.Interval {
		logger.Warn("The scenario did not run enough time to correct metrics. The metrics below might not be what you expect. Scenario Duration: %d, Report Interval: %d", gp.Duration, report.Interval)
	}
	report.PrintAllMetrics()
	inputs := scenarios.CollectInputParameters(gp.InputKeysForReport, gp.Raw.CommonParams, gp.Raw.ScenarioParams, gp.Raw.ParamsFromAPI)
	scenarioName := strings.Join([]string{ss.ScenarioName, ss.ScenarioPattern}, "-")
	report.WriteCSV(scenarioName, inputs)
	return nil
}
func run() error {
	err := setup()
	if err != nil {
		logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
		return err
	}
	err = start()
	if err != nil {
		logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
		return err
	}
	return nil
}
func main() {
	ss = &ScenarioSettings{HowMany: -1, Interval: -1, Duration: -1}
	isServerMode := util.GetEnv("BOT_SERVER_MODE")
	configPath := util.GetEnv("BOT_CONFIG")
	if configPath != "" {
		config = configPath
	}
	if isServerMode == "true" {
		err := listen()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}
	} else {
		err := load()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}
		err = run()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}
	}
}
