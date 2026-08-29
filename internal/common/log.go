package common

import (
	"fmt"

	"go.uber.org/zap"
)

const banner = "\n" +
	"								   ,-.----.                           \n" +
	"    ,---,.                         \\    /  \\      ,---,        ,---,. \n" +
	"  ,'  .' |                         |   :    \\   .'  .' `\\    ,'  .' | \n" +
	",---.'   |         ,--,      ,---, |   |  .\\ :,---.'     \\ ,---.'   | \n" +
	"|   |   .'       ,'_ /|  ,-+-. /  |.   :  |: ||   |  .`\\  ||   |   .' \n" +
	":   :  :    .--. |  | : ,--.'|'   ||   |   \\ ::   : |  '  |:   |  :   \n" +
	":   |  |-,,'_ /| :  . ||   |  ,\"' ||   : .   /|   | '  ;  ::   |  |-, \n" +
	"|   :  ;/||  ' | |  . .|   | /  | |;   | |`-' '   | ;  .  ||   :  ;/| \n" +
	"|   |   .'|  | ' |  | ||   | |  | ||   | ;    |   | :  |  '|   |   .' \n" +
	"'   :  '  :  | : ;  ; ||   | |  |/ :   | |    '   : | /  ; '   :  '   \n" +
	"|   |  |  '  :  `--'   \\   | |--'  :   : :    |   | '` ,/  |   |  |   \n" +
	"|   :  \\  :  ,      .-./   |/      |   | :    ;   :  .'    |   :  \\   \n" +
	"|   | ,'   `--`----'   '---'       `---'.|    |   ,.'      |   | ,'   \n" +
	"`----'                               `---`    '---'        `----'       \n"

var (
	logger *zap.Logger
)

func SyncLog() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not initialized")
	}
	logger.Fatal(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Warn(msg, fields...)
}

// Banner prints the ASCII art logo to stdout.
func Banner() {
	fmt.Print(banner)
}
