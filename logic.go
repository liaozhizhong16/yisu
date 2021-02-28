package yisu

/**
* @Time    : 21/2/4 下午1:45
* @Author  : liaozz
*
* 自我介绍一下
 */

import (
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

/**

const APP_URL = "https://www.baidu.com"  //常量
//如果是枚举类型的常量，需要先创建相应类型：
type Scheme string
const (
	HTTP  Scheme = "http"
	HTTPS Scheme = "https"
)
*/

type Data struct {
	Domain      string `json:"domain"` //请求获取RequestHost domain
	Time        string `json:"create_time"`
	RequestHost string `json:"request_host"` //数据存放地址
}
type ResponseData struct {
	Result     string `json:"result"`
	Desc       string `json:"desc"`
	ResultData Data   `json:"data"`
}
type MyLogic struct {
	Domain      string
	RequestHost string
	ConfigPath  string
	UUID        string
	debug       log.Logger
	info        log.Logger
	warn        log.Logger
}

func NewMyLocic(config_path string, registry prometheus.Gatherer, gather func(prometheus.Gatherer) string) *MyLogic {
	var read_config_t int16
	mylogic := &MyLogic{ConfigPath: config_path}
	mylogic.InitLog() //初始化log
	uuid, _ := get_uuid()
	mylogic.UUID = uuid
	mylogic.info.Log("uuid", uuid)
LOOP:
	mylogic.Init() //初始化
	read_config_t = 0

	for true {
		mylogic.PostProme(registry, gather)
		time.Sleep(time.Duration(5) * time.Second)
		read_config_t += 5
		if read_config_t > 10 {
			goto LOOP
		}
	}
	return mylogic
}

func (logic *MyLogic) PostProme(registry prometheus.Gatherer, gather func(prometheus.Gatherer) string) {
	if len(logic.RequestHost) > 0 {
		//收集数据请求到缓存服务器
		gather_data := gather(registry)
		params := url.Values{}
		params.Set("uuid", logic.UUID)
		params.Set("data", gather_data)
		params.Set("js", "exporter")
		//params= url.Values{"key": {"Value"}, "id": {"123"}}
		resp, err := http.PostForm(logic.RequestHost, params)
		if err != nil {
			logic.debug.Log("PostProme", err)
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			logic.debug.Log("PostProme", body)
		}

		logic.debug.Log("PostProme", gather_data)
	}
}
func (l *MyLogic) InitLog() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowAll()) // <--
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	l.debug = level.Debug(logger)
	l.info = level.Info(logger)
	l.warn = level.Warn(logger)
}
func (logic *MyLogic) CheckConfigPath() error {
	//判断配置文件是否存在
	if len(logic.ConfigPath) > 0 {
		if Exists(logic.ConfigPath) {
			return nil
		} else {
			logic.info.Log("CheckConfigPath", "config_path file not found")
			return errors.New("config_path file not found")
		}
	} else {
		logic.info.Log("CheckConfigPath", "config_path error")
		return errors.New("config_path error")
	}
}

func (logic *MyLogic) Init() {
	//读取配置文件
	//判断已读取就不再读取配置文件
	if len(logic.Domain) > 0 {
		logic.debug.Log("Init", "Domain="+logic.Domain)
	} else {
		if err := logic.CheckConfigPath(); err == nil { //检查配置文件
			//读取配置文件
			logic.ReadConfig()
		} else {
			logic.debug.Log("Init", err)
		}
	}
	//检查url地址是否有更改
	logic.CheckRequesUrl()
}

func (logic *MyLogic) CheckRequesUrl() {
	//请求检查地址
	if len(logic.Domain) > 0 {
		url := logic.Domain + "/assets/exporter/"
		resp, err := http.Get(url)
		if err != nil {
			logic.debug.Log("CheckRequesUrl", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				body, _ := ioutil.ReadAll(resp.Body)
				resd := &ResponseData{}
				logic.debug.Log("CheckRequesUrl", string(body))
				if err := json.Unmarshal(body, resd); err == nil {
					logic.debug.Log("CheckRequesUrl", "RequestHost="+resd.ResultData.RequestHost)
					logic.RequestHost = resd.ResultData.RequestHost

				} else {
					logic.info.Log("CheckRequesUrl", err)
				}
			} else {
				logic.info.Log("CheckRequesUrl", "request fail  resp.StatusCode="+string(resp.StatusCode))
			}
		}
	} else {
		logic.debug.Log("CheckRequesUrl", "Domain len is 0")
	}

}

func (logic *MyLogic) GetPromeDomain() (string, error) {
	//获取信息推送地址
	return "", nil
}

func (logic *MyLogic) ResponsePromeData() {
	//推动数据
}

func (logic *MyLogic) ReadConfig() {
	filePtr, err := os.Open(logic.ConfigPath)
	if err != nil {
		//panic(err)
		logic.info.Log("ReadConfig", err)
	} else {
		defer filePtr.Close()
		// 创建json解码器
		decoder := json.NewDecoder(filePtr)
		config_data := &Data{}
		err = decoder.Decode(config_data)
		if err != nil {
			logic.info.Log("ReadConfig", err)
		} else {
			logic.Domain = config_data.Domain
			logic.RequestHost = config_data.RequestHost
			logic.debug.Log("ReadConfig", "Domain="+config_data.Domain)
			logic.debug.Log("ReadConfig", "Time="+config_data.Time)
			logic.debug.Log("ReadConfig", "RequestHost="+config_data.RequestHost)
		}
	}

}
