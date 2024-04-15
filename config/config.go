package config

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config config

type config struct {
	DB           *DBConfig
	Graylog      *Graylog
	ScheduleCron string
}

func InitConfigDefault() {
	//logrus.Info("config initViper")
	InitViper("config", loadConfig)
	loadConfig()
}

func loadConfig() {

	// DBConfig
	Config.DB = &DBConfig{}
	Config.DB.Load("")

	// Graylog
	Config.Graylog = &Graylog{}
	Config.Graylog.Load("")

	Config.ScheduleCron = viper.GetString("system.scheduleCron")
}

// TCP TCP配置
type TCP struct {
	Port   int
	Enable bool
}

// Load 加载配置
func (t *TCP) Load() {
	t.Port = viper.GetInt("tcp.port")
	t.Enable = viper.GetBool("tcp.enable")
}

type DBConfig struct {
	DriverName   string
	Enable       bool
	Database     string
	User         string
	Password     string
	Charset      string
	Host         string
	Port         int
	ShowSql      bool
	MaxIdleConns int
	MaxOpenConns int
}

// URL mysql连接字符串
func (m *DBConfig) URL() string {

	if m.DriverName == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			m.User,
			m.Password,
			m.Host,
			m.Port,
			m.Database,
			m.Charset,
		)
	} else if m.DriverName == "postgres" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			m.Host,
			m.Port,
			m.User,
			m.Password,
			m.Database,
		)
	} else {
		return ""
	}
}

// Load 加载MySQL配置
func (m *DBConfig) Load(prefix string) {
	if prefix == "" {
		if viper.Get("mysql") == nil {
			prefix = "db"
		} else {
			prefix = "mysql"
		}
	}
	m.Enable = viper.GetBool(fmt.Sprintf("%s.enable", prefix))
	m.DriverName = viper.GetString(fmt.Sprintf("%s.driverName", prefix))
	m.Database = viper.GetString(fmt.Sprintf("%s.database", prefix))
	m.User = viper.GetString(fmt.Sprintf("%s.user", prefix))
	m.Password = viper.GetString(fmt.Sprintf("%s.password", prefix))
	m.Charset = viper.GetString(fmt.Sprintf("%s.charset", prefix))
	m.Host = viper.GetString(fmt.Sprintf("%s.host", prefix))
	m.Port = viper.GetInt(fmt.Sprintf("%s.port", prefix))
	m.MaxIdleConns = viper.GetInt(fmt.Sprintf("%s.maxIdleConns", prefix))
	m.MaxOpenConns = viper.GetInt(fmt.Sprintf("%s.maxOpenConns", prefix))
	m.ShowSql = viper.GetBool(fmt.Sprintf("%s.showSql", prefix))
}

// Redis redis配置
type Redis struct {
	Host     string
	Port     int
	Password string
	DBNum    int
	Idlesec  int
}

// URL redis连接字符串
func (r *Redis) URL() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// Load 加载Redis配置
func (r *Redis) Load(prefix string) {
	if prefix == "" {
		prefix = "redis"
	}
	r.Host = viper.GetString(fmt.Sprintf("%s.host", prefix))
	r.Port = viper.GetInt(fmt.Sprintf("%s.port", prefix))
	r.Password = viper.GetString(fmt.Sprintf("%s.password", prefix))
	r.DBNum = viper.GetInt(fmt.Sprintf("%s.dbNum", prefix))
	r.Idlesec = viper.GetInt(fmt.Sprintf("%s.idlesec", prefix))
}

// Graylog graylog配置
type Graylog struct {
	Host string
	Port int
}

// URL graylog连接字符换
func (g *Graylog) URL() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

// Load 加载Graylog配置
func (g *Graylog) Load(prefix string) {
	if prefix == "" {
		prefix = "graylog"
	}
	g.Host = viper.GetString(fmt.Sprintf("%s.host", prefix))
	g.Port = viper.GetInt(fmt.Sprintf("%s.port", prefix))
}

type Eterm3 struct {
	ServerIp   string
	Username   string
	Password   string
	CurIp      string
	Version    string
	VerifyCode string
	ServerPort int
}

func (e *Eterm3) Load(root, name string) {
	e.ServerIp = viper.GetString(fmt.Sprintf("%s.%s.serverIp", root, name))
	e.ServerPort = viper.GetInt(fmt.Sprintf("%s.%s.serverPort", root, name))
	e.Username = viper.GetString(fmt.Sprintf("%s.%s.username", root, name))
	e.Password = viper.GetString(fmt.Sprintf("%s.%s.password", root, name))
	e.Version = viper.GetString(fmt.Sprintf("%s.%s.version", root, name))
	e.VerifyCode = viper.GetString(fmt.Sprintf("%s.%s.verifyCode", root, name))
}

// InitViper 初始化配置
func InitViper(filename string, fn func()) {
	logrus.Debugf("开始初始化配置文件:%s", filename)
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Errorf("get dir err:%v", err.Error())
	}
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(currentDir)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("read config err:%v", err.Error())
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Errorf("本地配置更新:%s", e.Name)
		if fn != nil {
			fn()
		}
	})
}
