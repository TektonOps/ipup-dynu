package ip

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"gitub.com/khaliq/ddns/internal/config"
)

// var servers = []string{
// 	"https://api4.ipify.org",
// 	"https://ident.me",
// 	"https://ip2location.io/ip",
// }

// GetPublicIP tries to get the public IP address by querying a list of servers.
func GetPublicIP(logger *slog.Logger, conf *config.Config) (string, error) {
	for _, server := range conf.IPServers {
		resp, err := http.Get(server)
		if err != nil {
			logger.Warn("Failed to get IP from server", "server", server, "error", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Error reading response body", "server", server, "error", err)
				continue
			}

			ip := strings.TrimSpace(string(body))
			logger.Debug("Retrieved public IP", "ip", ip, "server", server)
			return ip, nil
		}

		logger.Warn("Server responded with non-200 status", "server", server, "status_code", resp.StatusCode)
	}

	return "", errors.New("could not retrieve public IP from any server")
}

// GetDomainIP resolves the IP address of the given domain name
func GetDomainIP(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	// Return the first IP address found
	if len(ips) > 0 {
		return ips[0].String(), nil
	}

	return "", nil
}
