package logs

import (
	"fmt"
	"log"
	"os"
)

var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var Info func(v ...interface{})
var Warn func(v ...interface{})
var Error func(v ...interface{})

func Initialize() error {
	infoLogger = log.New(log.Writer(), fmt.Sprintf("%s[INFO]%s: ", Green, Reset), log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(log.Writer(), fmt.Sprintf("%s[WARN]%s: ", Yellow, Reset), log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(log.Writer(), fmt.Sprintf("%s[ERROR]%s: ", Red, Reset), log.Ldate|log.Ltime|log.Lshortfile)

	if os.Getenv("DEBUG") == "true" {
		infoLogger.SetFlags(log.Ltime)
		warnLogger.SetFlags(log.Ltime)
		errorLogger.SetFlags(log.Ltime | log.Lshortfile)
		// set log output to os.Stdout
		infoLogger.SetOutput(os.Stdout)
		warnLogger.SetOutput(os.Stdout)
		errorLogger.SetOutput(os.Stdout)

	} else {
		// change prefix
		infoLogger.SetPrefix("[INFO]: ")
		warnLogger.SetPrefix("[WARN]: ")
		errorLogger.SetPrefix("[ERROR]: ")

		// set log flags
		infoLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		warnLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		errorLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		// create log files if not exists "logs/info.log", "logs/warn.log", "logs/error.log"
		infoLogFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error creating info.log file")
			return err
		}
		warnLogFile, err := os.OpenFile("logs/warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error creating warn.log file")
			return err
		}
		errorLogFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error creating error.log file")
			return err
		}

		// set log output to log files
		infoLogger.SetOutput(infoLogFile)
		warnLogger.SetOutput(warnLogFile)
		errorLogger.SetOutput(errorLogFile)
	}

	Info = infoLogger.Println
	Warn = warnLogger.Println
	Error = errorLogger.Println

	infoLogger.Println("Initialized loggers")
	return nil
}
