package collector

import (
	"os"
	"os/user"
	"time"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
)

// CollectHost gathers hostname, current user, and timestamp.
func CollectHost() (models.Host, error) {
	host := models.Host{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	hostname, err := os.Hostname()
	if err != nil {
		return host, err
	}
	host.Hostname = hostname

	if u, uerr := user.Current(); uerr == nil {
		host.CurrentUser = u.Username
	}

	return host, nil
}
