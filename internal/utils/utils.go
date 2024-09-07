package utils

import (
	"fmt"
	"log/slog"
	"strings"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func SetLogLevel(levelStr string) slog.Level {
	if level, ok := logLevelMap[strings.ToLower(levelStr)]; ok {
		return level
	}
	return slog.LevelInfo
}

func AppInfo(version, buildDate, commit, logLevel string) {
	cyan := "\x1b[36m"
	reset := "\x1b[0m"

	banner := `
━━┏━━━┓┏┓━┏┓━━━━
━━┃┏━┓┃┃┃━┃┃━━━━
┏┓┃┗━┛┃┃┃━┃┃┏━━┓
┣┫┃┏━━┛┃┃━┃┃┃┏┓┃
┃┃┃┃━━━┃┗━┛┃┃┗┛┃
┗┛┗┛━━━┗━━━┛┃┏━┛
━━━━━━━━━━━━┃┃━━
━━━━━━━━━━━━┗┛━━
  

	`
	fmt.Println(cyan + banner + reset)

	fmt.Printf("Version: %v\n\nBuild Date: %v\n\nCommit: %v\n\n", cyan+version+reset, cyan+buildDate+reset, cyan+commit+reset)

	fmt.Printf("Log Level Set to: "+cyan+"%v\n\n"+reset, logLevel)

}
