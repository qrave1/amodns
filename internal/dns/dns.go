package dns

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	stageDNS = "10.13.245.31"
	prodDNS  = "10.13.245.30"
)

type DNSchanger struct{}

func NewDNSchanger() *DNSchanger {
	return &DNSchanger{}
}

func (d *DNSchanger) MapEnvToAddr(env string) string {
	switch env {
	case "stage":
		return stageDNS
	case "prod":
		return prodDNS
	default:
		return ""
	}
}

// GetActiveConnectionName получает активное соединение Wi-Fi (интерфейс)
func (d *DNSchanger) GetActiveConnectionName() (string, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "DEVICE,TYPE,STATE", "device", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var interfaceName string

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 3 && fields[1] == "wifi" && fields[2] == "connected" {
			interfaceName = fields[0]
		}
	}

	cmd = exec.Command("nmcli", "-t", "-f", "NAME,DEVICE", "connection", "show", "--active")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 2 && fields[1] == interfaceName {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("connection name for %s not found", interfaceName)
}

// SetDNS изменяет DNS сервера для указанного соединения
func (d *DNSchanger) SetDNS(connectionName string, dnsServer string) error {
	cmd := exec.Command("nmcli", "connection", "modify", connectionName, "ipv4.dns", dnsServer)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error exec command: %v\nOutput: %s", err, output)
	}

	cmd = exec.Command("nmcli", "connection", "up", connectionName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error exec command: %v\nOutput: %s", err, output)
	}

	return nil
}

// GetCurrentDNS получает текущие DNS-серверы для активного соединения.
func (d *DNSchanger) GetCurrentDNS() (string, error) {
	cmd := exec.Command("nmcli", "device", "show")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения команды nmcli: %v", err)
	}

	// Парсим вывод, чтобы найти строки с "IP4.DNS"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "IP4.DNS") {
			// Разделяем строку по пробелу и берём последний элемент (адрес DNS-сервера)
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("error get current DNS")
}
