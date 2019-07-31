package init

import (
	"ehelp/setting"
	"ehelp/x/db/mongodb"
	"ehelp/x/fcm"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"

	"io"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	load()
}

func load() {
	// Open our jsonFile
	jsonFile, err := os.Open("app.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result configApp
	json.Unmarshal([]byte(byteValue), &result)
	//initLog(result.Log)

	initDB(result.DB)
	initFcm(result.Fcm)
	initSystem(result.System)
}

type configApp struct {
	DB     db           `json:"db"`
	Fcm    fcmConfig    `json:"fcm"`
	System systemConfig `json:"system"`
	Log    logConfig    `json:"log"`
}
type db struct {
	Path    string `json:"path"`
	DBName  string `json:"db_name"`
	MaxPool int    `json:"max_pool"`
}

type logConfig struct {
	LogDir          string `json:"log_dir"`
	Alsologtostderr string `json:"alsologtostderr"`
}
type fcmConfig struct {
	FcmEmp string `json:"fcm_employee"`
	FcmCus string `json:"fcm_customer"`
	Owner  string `json:"owner"`
	Avatar string `json:"avatar"`
}
type systemConfig struct {
	TimeHourHiddenOrder   int     `json:"time_hour_hidden_order"`
	AboutHourGoWork       int     `json:"about_hour_go_work"`
	AboutMinuteFinishWork float64 `json:"about_minute_finish_work"`
	AboutMinuteWorking    float64 `json:"about_minute_working"`
}

func initLog(l logConfig) {

	//config for gin request log
	{
		f, _ := os.Create(path.Join("log", "gin.log"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		defer f.Close()
	}
	//config for app log use glog
	{
		flag.Set("alsologtostderr", l.Alsologtostderr)
		flag.Set("log_dir", l.LogDir)
		flag.Parse()
	}
	if _, err := os.Stat(l.LogDir); os.IsNotExist(err) {
		os.Mkdir(l.LogDir, os.ModeAppend)
	}
	flag.Parse()
}

func initDB(database db) {
	fmt.Println(database)
	// Read configuration.
	// mongodb.MaxPool = context.IntDefault("mongo.maxPool", 0)
	// mongodb.PATH, _ = context.String("mongo.path")
	// mongodb.DBNAME, _ = context.String("mongo.database")
	mongodb.MaxPool = database.MaxPool
	mongodb.PATH = database.Path
	mongodb.DBNAME = database.DBName
	mongodb.CheckAndInitServiceConnection()
}

func initFcm(f fcmConfig) {
	// fcm.FCM_SERVER_KEY_CUSTOMER, _ = context.String("fcm.serverkey.customer")
	// fcm.FCM_SERVER_KEY_EMPLOYEE, _ = context.String("fcm.serverkey.employee")
	// fcm.LINK_AVATAR, _ = context.String("server.avatar")
	fcm.FCM_SERVER_KEY_CUSTOMER = f.FcmEmp
	fcm.FCM_SERVER_KEY_EMPLOYEE = f.FcmCus
	fcm.LINK_AVATAR = f.Avatar
	fcm.NewFcmApp(fcm.FCM_SERVER_KEY_CUSTOMER, fcm.FCM_SERVER_KEY_EMPLOYEE)
}

func initSystem(s systemConfig) {
	// setting.SettingSys.TimeHourHiddenOrder, _ = context.Int("server.time_hour_hidden_order")
	// setting.SettingSys.AboutHourGoWork, _ = context.Int("server.about_hour_go_work")
	// var finish, _ = context.String("server.about_minute_finish_work")
	// setting.SettingSys.AboutMinuteFinishWork, _ = strconv.ParseFloat(finish, 64)
	// var working, _ = context.String("server.about_minute_working")
	// setting.SettingSys.AboutMinuteWorking, _ = strconv.ParseFloat(working, 64)
	setting.SettingSys.TimeHourHiddenOrder = s.TimeHourHiddenOrder
	setting.SettingSys.AboutMinuteWorking = s.AboutMinuteWorking
	setting.SettingSys.AboutMinuteFinishWork = s.AboutMinuteFinishWork
	setting.SettingSys.AboutHourGoWork = s.AboutHourGoWork

}

// func initCache() {
// 	rest.AssertNil(cache.SetCacheSystem())
// }
