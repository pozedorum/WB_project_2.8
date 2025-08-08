package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	fileLogger *log.Logger // Для записи в файл
	consoleLog bool        // Дублировать в консоль

}

func New(logFile string, consoleLog bool) (*Logger, error) {

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		fileLogger: log.New(file, "", log.LstdFlags),
		consoleLog: consoleLog,
	}, nil
}

// Middleware для HTTP-запросов
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		l.log("[REQUEST]", r.Method, r.URL.Path, r.RemoteAddr, duration)

	})
}

func (l *Logger) LogError(r *http.Request, err error) {
	l.log("[ERROR]", r.Method, r.URL.Path, err.Error())

}

func (l *Logger) log(prefix string, v ...interface{}) {
	msg := prefix + " " + fmt.Sprint(v...)
	l.fileLogger.Println(msg)

	if l.consoleLog {
		log.Println(msg)
	}
}
