package collector

import (
	"net"
	"strings"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"

	psnet "github.com/shirou/gopsutil/v3/net"
)

// CollectNetwork gathers network interface details.
func CollectNetwork() ([]models.NetworkInterface, error) {
	var result []models.NetworkInterface

	ifaces, err := psnet.Interfaces()
	if err != nil {
		return result, err
	}

	for _, iface := range ifaces {
		// Skip loopback interfaces
		if isLoopback(iface.Flags) {
			continue
		}

		ni := models.NetworkInterface{
			Name:       iface.Name,
			MACAddress: iface.HardwareAddr,
			Type:       classifyInterface(iface.Name),
			IPv4:       []string{},
			IPv6:       []string{},
		}

		for _, addr := range iface.Addrs {
			ip := stripCIDR(addr.Addr)
			parsed := net.ParseIP(ip)
			if parsed == nil {
				continue
			}
			if parsed.To4() != nil {
				ni.IPv4 = append(ni.IPv4, ip)
			} else {
				ni.IPv6 = append(ni.IPv6, ip)
			}
		}

		result = append(result, ni)
	}

	return result, nil
}

func isLoopback(flags []string) bool {
	for _, f := range flags {
		if strings.EqualFold(f, "loopback") {
			return true
		}
	}
	return false
}

func stripCIDR(addr string) string {
	if idx := strings.Index(addr, "/"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// classifyInterface makes a best-effort guess at the interface type
// based on common naming conventions across platforms.
func classifyInterface(name string) string {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "wi-fi") || strings.Contains(n, "wlan") ||
		strings.Contains(n, "wifi") || strings.HasPrefix(n, "wl"):
		return "WiFi"
	case strings.Contains(n, "vpn") || strings.Contains(n, "tun") ||
		strings.Contains(n, "tap") || strings.Contains(n, "ppp") ||
		strings.Contains(n, "wg"):
		return "VPN"
	case strings.Contains(n, "veth") || strings.Contains(n, "docker") ||
		strings.Contains(n, "br-") || strings.Contains(n, "virbr") ||
		strings.Contains(n, "vmnet") || strings.Contains(n, "vbox"):
		return "Virtual"
	case strings.HasPrefix(n, "eth") || strings.HasPrefix(n, "en") ||
		strings.Contains(n, "ethernet"):
		return "Ethernet"
	default:
		return "Unknown"
	}
}
