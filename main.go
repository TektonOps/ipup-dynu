package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"gitub.com/khaliq/ddns/internal/config"
	"gitub.com/khaliq/ddns/internal/dynu"
	"gitub.com/khaliq/ddns/internal/ip"
	"gitub.com/khaliq/ddns/internal/utils"
)

var (
	version   string
	buildDate string
	commitSha string
)

func main() {
	conf := config.New()

	logLevel := utils.SetLogLevel(conf.LogLevel)
	opts := &tint.Options{
		Level:     logLevel,
		AddSource: conf.Logs.EnableSource,
	}
	logger := slog.New(tint.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	utils.AppInfo(version, buildDate, commitSha, conf.LogLevel)

	var lastIP string

	for {
		// Check the IP that the domain currently resolves to
		domainIP, err := ip.GetDomainIP(conf.DomainName)
		if err != nil {
			logger.Error("Failed to resolve domain IP", "error", err)
			time.Sleep(conf.CheckInterval * time.Second)
			continue
		}

		// Get the current public IP address
		currentIP, err := ip.GetPublicIP(logger, conf)
		if err != nil {
			logger.Error("Failed to get public IP", "error", err)
			time.Sleep(conf.CheckInterval * time.Second)
			continue
		}

		// If the IP has changed or is being checked for the first time, update DNS
		if lastIP != currentIP || domainIP != currentIP {
			req, err := dynu.NewDNSRequest(conf.UserName, conf.Password, conf.DomainName, conf.GroupName, currentIP, conf.EnableGroup, logger)
			if err != nil {
				logger.Error("Error initializing DNS request", "error", err)
				time.Sleep(conf.CheckInterval * time.Second)
				continue
			}

			response, err := dynu.UpdateDNSRecord(req, logger)
			if err != nil {
				logger.Error("Error updating DNS record", "error", err)

				// Exit the program if authentication fails
				if strings.Contains(err.Error(), "authentication failed") {
					os.Exit(1)
				}
			} else {
				logger.Info("DNS update status", "response", response)
				lastIP = currentIP
			}
		}

		time.Sleep(conf.CheckInterval * time.Second)
	}
}
