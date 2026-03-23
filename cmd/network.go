package cmd

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

var TempNetworkName string
var kvmageCreatedNetwork bool

func listLibvirtNetworks() ([]string, error) {
	out, err := exec.Command("virsh", "net-list", "--all", "--name").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list libvirt networks: %w", err)
	}
	var networks []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			networks = append(networks, line)
		}
	}
	return networks, nil
}

func listSystemBridges() []string {
	var bridges []string
	// Get bridges from network interfaces on the system
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			bridges = append(bridges, iface.Name)
		}
	}
	// Also get bridge names from libvirt network XML definitions
	networks, _ := listLibvirtNetworks()
	for _, name := range networks {
		out, err := exec.Command("virsh", "net-dumpxml", name).Output()
		if err != nil {
			continue
		}
		xml := string(out)
		for _, line := range strings.Split(xml, "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "<bridge name=") {
				nameStart := strings.Index(line, "name='")
				if nameStart == -1 {
					nameStart = strings.Index(line, "name=\"")
				}
				if nameStart != -1 {
					nameStart += 6
					nameEnd := strings.IndexAny(line[nameStart:], "'\"")
					if nameEnd != -1 {
						bridges = append(bridges, line[nameStart:nameStart+nameEnd])
					}
				}
			}
		}
	}
	return bridges
}

func usedSubnets() ([]string, error) {
	networks, err := listLibvirtNetworks()
	if err != nil {
		return nil, err
	}
	var subnets []string
	for _, name := range networks {
		out, err := exec.Command("virsh", "net-dumpxml", name).Output()
		if err != nil {
			continue
		}
		xml := string(out)
		for _, line := range strings.Split(xml, "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "<ip address=") {
				addrStart := strings.Index(line, "address='")
				if addrStart == -1 {
					addrStart = strings.Index(line, "address=\"")
				}
				if addrStart != -1 {
					addrStart += 9
					addrEnd := strings.IndexAny(line[addrStart:], "'\"")
					if addrEnd != -1 {
						subnets = append(subnets, line[addrStart:addrStart+addrEnd])
					}
				}
			}
		}
	}
	return subnets, nil
}

func findAvailableBridge(existing []string) string {
	for i := 0; ; i++ {
		candidate := fmt.Sprintf("virbr%d", i)
		found := false
		for _, b := range existing {
			if b == candidate {
				found = true
				break
			}
		}
		if !found {
			return candidate
		}
	}
}

func findAvailableSubnet(usedAddrs []string) (string, string, string) {
	for i := 122; i <= 254; i++ {
		gateway := fmt.Sprintf("192.168.%d.1", i)
		conflict := false
		for _, addr := range usedAddrs {
			if strings.HasPrefix(addr, fmt.Sprintf("192.168.%d.", i)) {
				conflict = true
				break
			}
		}
		if !conflict {
			rangeStart := fmt.Sprintf("192.168.%d.2", i)
			rangeEnd := fmt.Sprintf("192.168.%d.254", i)
			return gateway, rangeStart, rangeEnd
		}
	}
	// Fallback to 10.x range
	for i := 200; i <= 254; i++ {
		gateway := fmt.Sprintf("10.%d.%d.1", i, i)
		conflict := false
		for _, addr := range usedAddrs {
			if strings.HasPrefix(addr, fmt.Sprintf("10.%d.%d.", i, i)) {
				conflict = true
				break
			}
		}
		if !conflict {
			rangeStart := fmt.Sprintf("10.%d.%d.2", i, i)
			rangeEnd := fmt.Sprintf("10.%d.%d.254", i, i)
			return gateway, rangeStart, rangeEnd
		}
	}
	return "", "", ""
}

func EnsureKvmageNetwork() (string, error) {
	networks, err := listLibvirtNetworks()
	if err != nil {
		return "", err
	}

	// Check if an existing kvmage network is already active — reuse it
	for _, n := range networks {
		if strings.HasPrefix(n, "kvmage-") {
			out, err := exec.Command("virsh", "net-info", n).Output()
			if err == nil && strings.Contains(string(out), "Active:         yes") {
				PrintVerbose(2, "Using existing kvmage network: %s", n)
				TempNetworkName = n
				return n, nil
			}
		}
	}

	// No active kvmage network — create one using the same random string as TempImageName
	netName := fmt.Sprintf("kvmage-%s", strings.TrimPrefix(TempImageName, "kvmage-"))

	bridges := listSystemBridges()
	subnets, err := usedSubnets()
	if err != nil {
		return "", err
	}

	bridgeName := findAvailableBridge(bridges)
	gateway, rangeStart, rangeEnd := findAvailableSubnet(subnets)

	if gateway == "" {
		return "", fmt.Errorf("could not find an available subnet for kvmage network")
	}

	xml := fmt.Sprintf(`<network>
  <name>%s</name>
  <bridge name='%s'/>
  <forward mode='nat'/>
  <dns enable='yes'/>
  <ip address='%s' netmask='255.255.255.0'>
    <dhcp>
      <range start='%s' end='%s'/>
    </dhcp>
  </ip>
</network>`, netName, bridgeName, gateway, rangeStart, rangeEnd)

	PrintVerbose(2, "Creating kvmage network: %s (bridge: %s, subnet: %s/24)", netName, bridgeName, gateway)

	defineCmd := exec.Command("virsh", "net-define", "/dev/stdin")
	defineCmd.Stdin = strings.NewReader(xml)
	if out, err := defineCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to define network %s: %s: %w", netName, string(out), err)
	}

	if out, err := exec.Command("virsh", "net-start", netName).CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to start network %s: %s: %w", netName, string(out), err)
	}

	PrintVerbose(1, "Created kvmage network: %s", netName)
	TempNetworkName = netName
	kvmageCreatedNetwork = true
	return netName, nil
}

func CleanupKvmageNetwork() {
	// Only clean up networks that kvmage created this session
	if !kvmageCreatedNetwork || TempNetworkName == "" {
		return
	}

	// Check if any kvmage VMs are still using this network
	out, err := exec.Command("virsh", "list", "--name").Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			name := strings.TrimSpace(line)
			if strings.HasPrefix(name, "kvmage-") {
				PrintVerbose(2, "Skipping network cleanup: VM %s is still running", name)
				return
			}
		}
	}

	PrintVerbose(2, "Destroying kvmage network: %s", TempNetworkName)
	exec.Command("virsh", "net-destroy", TempNetworkName).Run()
	exec.Command("virsh", "net-undefine", TempNetworkName).Run()
	TempNetworkName = ""
	kvmageCreatedNetwork = false
}

func cleanupOrphanedNetworks() {
	networks, err := listLibvirtNetworks()
	if err != nil {
		return
	}
	for _, name := range networks {
		if !strings.HasPrefix(name, "kvmage-") {
			continue
		}
		PrintVerbose(2, "Removing orphaned kvmage network: %s", name)
		exec.Command("virsh", "net-destroy", name).Run()
		exec.Command("virsh", "net-undefine", name).Run()
		Print("Removed network: %s", name)
	}
}
