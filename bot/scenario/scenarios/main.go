package scenarios

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Diarkis/diarkis/util"
	uuidv4 "github.com/Diarkis/diarkis/uuid/v4"
	"handson/bot/scenario/lib/log"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type UserState struct {
	sync.RWMutex
	state map[string]map[string]any
}
type GlobalParams struct {
	ScenarioName       string         `json:"scenarioName"`
	ScenarioPattern    string         `json:"scenarioPattern"`
	Interval           int            `json:"interval"`
	Duration           int            `json:"duration"`
	IdleDuration       int            `json:"idleDuration"`
	HowMany            int            `json:"howmany"`
	Host               string         `json:"host"`
	LogLevel           int            `json:"logLevel"`
	ReceiveByteSize    int            `json:"receiveByteSize"`
	UDPSendInterval    int64          `json:"udpSendInterval"`
	MetricsInterval    int            `json:"metricsInterval"`
	Configs            map[string]int `json:"configs"`
	InputKeysForReport []string       `json:"keysForReport"`
	Raw                struct {
		CommonParams   map[string]any
		ScenarioParams map[string]any
		ParamsFromAPI  map[string]any
	}
	sync.RWMutex
	UserState *UserState
}
type Scenario interface {
	GetUserID() string
	Run(globalParams *GlobalParams) error
	ParseParam(index int, params []byte) error
	OnScenarioEnd() error
	OnIdle()
}

var logger = log.New("BOT/SCENARIO")
var ScenarioFactoryList map[string]func() Scenario = map[string]func() Scenario{"Connect": NewConnectScenario, "Ticket": NewTicketScenario, "Session": NewSessionScenario}

type ApiParamAttributes struct {
	DefaultValue any   `json:"defaultValue"`
	Options      []any `json:"options"`
	OptionRates  []any `json:"optionRates"`
	Range        struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"range"`
	IsRandom     bool   `json:"isRandom"`
	IsSequential bool   `json:"isSequential"`
	Type         string `json:"type"`
}

func (ap *ApiParamAttributes) isRangeSet() bool {
	return ap.Range.Min != 0 && ap.Range.Max != 0
}
func NewUserState() *UserState {
	us := &UserState{}
	us.state = map[string]map[string]any{}
	return us
}
func (us *UserState) Init(userID string) {
	us.Lock()
	defer us.Unlock()
	us.state[userID] = map[string]any{}
}
func (us *UserState) Get(userID string, key string) any {
	us.RLock()
	defer us.RUnlock()
	if itf, ok := us.state[userID]; ok {
		if val, ok := itf[key]; ok {
			return val
		}
	}
	return nil
}
func (us *UserState) Set(userID string, key string, value any) error {
	us.Lock()
	defer us.Unlock()
	if _, ok := us.state[userID]; !ok {
		return errors.New("User does not exist. Call Init first.")
	}
	us.state[userID][key] = value
	return nil
}
func (us *UserState) Search(key string, value any, limit int) []string {
	us.RLock()
	defer us.RUnlock()
	var list []string
	for userID, state := range us.state {
		if v, _ := state[key]; v == value {
			list = append(list, userID)
			if len(list) == limit {
				break
			}
		}
	}
	return list
}
func (us *UserState) Exists(userID string) bool {
	us.RLock()
	defer us.RUnlock()
	_, ok := us.state[userID]
	return ok
}
func (gp *GlobalParams) GenerateParams(index int) ([]byte, error) {
	return GenerateParams(index, gp.Raw.CommonParams, gp.Raw.ScenarioParams, gp.Raw.ParamsFromAPI)
}
func GenerateParams(index int, params ...map[string]any) ([]byte, error) {
	ret := map[string]any{}
	uuid, _ := uuidv4.New()
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	for _, paramsMap := range params {
		for key, paramIf := range paramsMap {
			if param, ok := paramIf.(map[string]any); !ok {
				ret[key] = paramIf
			} else {
				apiParam := ApiParamAttributes{}
				jsonParams, err := json.Marshal(param)
				if err != nil {
					continue
				}
				err = json.Unmarshal(jsonParams, &apiParam)
				if err != nil {
					continue
				}
				options := apiParam.Options
				optionRates := apiParam.OptionRates
				if len(options) > 0 && len(optionRates) > 0 {
					randomValue := rand.Float64()
					c := 0.0
					for i, option := range options {
						var rate float64
						if i < len(optionRates) {
							if rate_, ok := optionRates[i].(float64); ok {
								rate = rate_
							}
						}
						c += rate
						if randomValue < c || i == len(options)-1 {
							ret[key] = option
							break
						}
					}
				} else if apiParam.IsRandom {
					if len(options) > 0 {
						ret[key] = options[util.RandomInt(0, len(options)-1)]
					} else {
						ret[key] = uuid.String
					}
				} else if apiParam.IsSequential {
					baseValue, ok := apiParam.DefaultValue.(float64)
					if ok {
						v := int(baseValue)
						if apiParam.isRangeSet() {
							delta := apiParam.Range.Max - apiParam.Range.Min + 1
							ret[key] = int(math.Mod(float64(v), float64(delta))) + apiParam.Range.Min
							continue
						}
						ret[key] = index + v
						continue
					}
					baseValueStr, ok := apiParam.DefaultValue.(string)
					if ok {
						baseValue, err := strconv.Atoi(baseValueStr)
						if err == nil {
							if apiParam.isRangeSet() {
								delta := apiParam.Range.Max - apiParam.Range.Min + 1
								ret[key] = int(math.Mod(float64(baseValue), float64(delta))) + apiParam.Range.Min
								continue
							}
							ret[key] = strconv.Itoa(index + baseValue)
							continue
						}
					}
					ret[key] = index
				} else if apiParam.isRangeSet() {
					ret[key] = util.RandomInt(apiParam.Range.Min, apiParam.Range.Max)
				} else {
					if apiParam.DefaultValue != nil {
						ret[key] = apiParam.DefaultValue
					} else {
						ret[key] = param
					}
				}
				if apiParam.Type != "" {
					ret[key] = DynamicCast(apiParam.Type, ret[key])
				}
			}
		}
	}
	bytes, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func CollectInputParameters(keys []string, params ...map[string]any) map[string]string {
	ret := map[string]string{}
	for _, input := range params {
		for _, key := range keys {
			if paramIf, ok := input[key]; ok {
				if param, ok := paramIf.(map[string]any); ok {
					apiParam := ApiParamAttributes{}
					jsonParams, err := json.Marshal(param)
					if err != nil {
						continue
					}
					err = json.Unmarshal(jsonParams, &apiParam)
					if err != nil {
						continue
					}
					if apiParam.IsRandom {
						ret[key] = "Random"
						continue
					}
					if apiParam.IsSequential {
						ret[key] = "Sequential"
						continue
					}
					if apiParam.isRangeSet() {
						ret[key] = fmt.Sprintf("Range [%d-%d]", apiParam.Range.Min, apiParam.Range.Max)
					}
					ret[key] = fmt.Sprintf("%v", paramIf)
				} else {
					ret[key] = fmt.Sprintf("%v", paramIf)
				}
			}
		}
	}
	return ret
}
func DynamicCast(sType string, target any) any {
	return target
}
