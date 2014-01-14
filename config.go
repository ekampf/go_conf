package go_conf

import (
	"flag"
	"github.com/kylelemons/go-gypsy/yaml"
	"io/ioutil"
	"log"
	"os"
  "path"
)

var (
  exitHandler   ExitHandler
  config        *yaml.File
  environment   string
  config_file   = flag.String("config", "./config/config.yml", "the config.yml")
  log_file_path = flag.String("log", "./log/", "where does the log go?")
)

var GlobalLogFile *os.File

func init() {
	exitHandler = &StandardHandler{}
	flag.Parse()
	setEnv()
	initlogAndConfig()
	startSignalCatcher()
}

func initlogAndConfig() {
	//create log
  log_file_name := GetLogFile()
	log_file, err := os.OpenFile(log_file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("cannot write log")
	}
	log.SetOutput(log_file)
	log.SetFlags(5)
  
  GlobalLogFile = log_file

	//read the config and build config stuff
	c_file, err := ioutil.ReadFile(*config_file)
	if err != nil {
		log.Panic("no config file found")
	}
	config = yaml.Config(string(c_file))
}

func GetLogFile() string {
  return path.Join(*log_file_path, environment + ".log")
}

func GetEnv() string {
	return environment
}

func setEnv() {
	environment = os.Getenv("GO_ENV")
	if environment == "" {
		environment = "development"
	}
}

func getConfigParameter(prefix, name string) string {
  yml_param := environment + "."
  if (prefix != "") {
    yml_param =  yml_param + prefix + "." + name
  } else {
    yml_param = yml_param + name
  }
  
  param, err := config.Get(yml_param)
  if err != nil {
  	log.Panic("missing config parameter: " + yml_param)
  }
  return param
}

func GetConfigParameter(prefix, name string) string {
  return getConfigParameter(prefix, name)
}
