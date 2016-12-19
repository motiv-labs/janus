package proxy

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Proxy defines proxy rules for a route
type Proxy struct {
	PreserveHostHeader          bool     `bson:"preserve_host_header" json:"preserve_host_header"`
	ListenPath                  string   `bson:"listen_path" json:"listen_path" valid:"required"`
	TargetURL                   string   `bson:"target_url" json:"target_url" valid:"url,required"`
	StripListenPath             bool     `bson:"strip_listen_path" json:"strip_listen_path"`
	EnableLoadBalancing         bool     `bson:"enable_load_balancing" json:"enable_load_balancing"`
	TargetList                  []string `bson:"target_list" json:"target_list"`
	CheckHostAgainstUptimeTests bool     `bson:"check_host_against_uptime_tests" json:"check_host_against_uptime_tests"`
	Methods                     []string `bson:"methods" json:"methods"`
}

//Validate validates proxy data
func Validate(proxy Proxy) bool {
	if proxy.ListenPath == "" {
		log.Warning("Listen path is empty")
		return false
	}

	if strings.Contains(proxy.ListenPath, " ") {
		log.Warning("Listen path contains spaces, is invalid")
		return false
	}

	return true
}
