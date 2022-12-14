package conf

import (
	"log"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
	DB       string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPasswd string
	DBName   string

	IsUseTestData  bool
	Number         int
	ProjectKey     string
	StartTimestamp int64
	EndTimestamp   int64
)

func Init() {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		log.Fatalln("加载配置文件失败，请检查")
	}
	log.Println("加载配置文件成功")
	loadServer(file)
	loadMysql(file)
	loadOther(file)
	dsn := DBUser + ":" + DBPasswd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println("正在连接数据库:", dsn)
	model.Database(dsn)
	log.SetFlags(log.Ldate | log.Llongfile | log.Ltime)
	if IsUseTestData {
		utils.GenerateTestData(Number, ProjectKey, StartTimestamp, EndTimestamp)
	}
}

func loadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func loadMysql(file *ini.File) {
	DB = file.Section("mysql").Key("DB").String()
	DBHost = file.Section("mysql").Key("DBHost").String()
	DBPort = file.Section("mysql").Key("DBPort").String()
	DBUser = file.Section("mysql").Key("DBUser").String()
	DBPasswd = file.Section("mysql").Key("DBPasswd").String()
	DBName = file.Section("mysql").Key("DBName").String()
}

func loadOther(file *ini.File) {
	other := file.Section("other")
	var err error
	IsUseTestData, err = other.Key("UseTestData").Bool()
	Number, err = other.Key("Number").Int()
	ProjectKey = other.Key("ProjectKey").String()
	StartTimestamp, err = other.Key("StartTimestamp").Int64()
	EndTimestamp, err = other.Key("EndTimestamp").Int64()
	if err != nil {
		log.Fatalln("读取Other配置项失败，请检查", err)
	}
}
