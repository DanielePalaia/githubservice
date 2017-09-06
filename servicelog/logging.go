package servicelog

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	filename string
	*log.Logger
}

var logger *Logger
var once sync.Once

// start loggeando
func GetInstance() *Logger {
	once.Do(func() {
		logger = createLogger("./logs/mylogger.log")
	})
	return logger
}

func createLogger(fname string) *Logger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &Logger{
		filename: fname,
		Logger:   log.New(file, "NewRelicService ", log.Lshortfile),
	}
}
