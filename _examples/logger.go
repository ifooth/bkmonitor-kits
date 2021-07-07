package main

import "github.com/TencentBlueKing/bkmonitor-kits/logger"

// InitLogger initializes the logger
func InitLogger() {
	// feel free to config the options.
	logger.SetOptions(logger.Options{
		Filename:   "/data/log/myproject/applog",
		MaxSize:    1000, // 1GB
		MaxAge:     3,    // 3 days
		MaxBackups: 3,    // 3 backups
	})
}

func main() {
	// Note: init logger when you want to run your program in the production env.
	// for example:
	InitLogger()

	// use logger Method anywhere directly, such as Info/Warn/Error/...
	// logs will be displayed on the stdout stream by default.
	logger.Info("This is the info level message.")
	logger.Warnf("This is the warn level message. %s", "oop!")
	logger.Error("Something error here.")
}
