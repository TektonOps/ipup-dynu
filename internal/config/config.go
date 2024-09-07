package config

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var conf Config

// Config represents the configuration structure
type Config struct {
	Dynu `yaml:"dynu"`
	Logs `yaml:"logs"`
}

type Dynu struct {
	DomainName    string        `yaml:"domain"`
	GroupName     string        `yaml:"group"`
	EnableGroup   bool          `yaml:"enableGroup"`
	UserName      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	IPServers     []string      `yaml:"ipServersList"`
	CheckInterval time.Duration `yaml:"ipCheckInterval"`
}

type Logs struct {
	LogLevel     string `yaml:"logLevel"`
	EnableSource bool   `yaml:"enableSource"`
}

// New creates a new Config instance based on environment variables or YAML file
func New() *Config {
	useConfig := envLookupBool("USE_CONFIG_FILE", false)
	if useConfig {
		return loadConfig()
	}

	conf.Dynu = Dynu{
		DomainName:    envLookup("DYNU_DOMAIN_NAME", "example.com"),
		GroupName:     envLookup("DYNU_GROUP_NAME", ""),
		UserName:      envLookup("DYNU_USERNAME", ""),
		Password:      envLookup("DYNU_PASSWORD", ""),
		IPServers:     envLookupIPServers("IPSERVERS_LIST", "https://api4.ipify.org,https://ip2location.io/ip,https://ident.me"),
		EnableGroup:   envLookupBool("DYNU_ENABLE_GROUP", false),
		CheckInterval: envLookupInterval("DYNU_IPCHECK_INTERVAL", "60"),
	}

	conf.Logs = Logs{
		LogLevel:     envLookup("LOG_LEVEL", "info"),
		EnableSource: envLookupBool("ENABLE_LOG_SOURCE", false),
	}

	return &conf
}

// envLookup reads a string value or returns the default
func envLookup(key string, defOption string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defOption
}

// envLookup reads a string value or returns the default
func envLookupIPServers(key string, defOption string) []string {
	if value, exists := os.LookupEnv(key); exists {
		valList := strings.Split(value, ",")
		return valList
	}
	defaultOptions := strings.Split(defOption, ",")
	return defaultOptions
}

// envLookupInterval reads a string value and  returns time.duration
func envLookupInterval(key string, defOption string) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		sleepDuration, _ := strconv.Atoi(value)

		return time.Duration(sleepDuration)
	}
	dv, _ := strconv.Atoi(defOption)

	return time.Duration(dv)
}

// envLookupBool reads a bool value or returns the default
func envLookupBool(name string, defOption bool) bool {
	valStr := envLookup(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defOption
}

// loadConfig opens and reads the "config.yaml" file and unmarshals
// its content into the Config struct.
func loadConfig() *Config {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %s", err)
	}
	return &conf
}
