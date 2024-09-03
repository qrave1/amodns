package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Использование: %s [stage|prod]", os.Args[0])
	}

	env := os.Args[1]
	var dnsServers []string

	switch env {
	case "stage":
		dnsServers = []string{"10.13.245.31"}
	case "prod":
		dnsServers = []string{"10.13.245.30"}
	default:
		log.Fatalf("Неизвестное окружение: %s. Используйте 'stage' или 'prod'.", env)
	}

	// Получаем текущее активное соединение Wi-Fi
	network, err := getActiveNetwork()
	if err != nil {
		log.Fatalf("Ошибка при получении активного соединения: %v", err)
	}

	fmt.Printf("Активное соединение Wi-Fi: %s\n", network)

	// Получаем имя соединения
	connectionName, err := getConnectionName(network)
	if err != nil {
		log.Fatalf("Ошибка при получении имени соединения: %v", err)
	}

	fmt.Printf("Имя соединения: %s\n", connectionName)

	// Меняем DNS на выбранный
	err = setDNS(connectionName, dnsServers)
	if err != nil {
		log.Fatalf("Ошибка при установке DNS: %v", err)
	}

	fmt.Println("DNS успешно изменен.")
}

// getActiveNetwork получает активное соединение Wi-Fi (интерфейс)
func getActiveNetwork() (string, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "DEVICE,TYPE,STATE", "device", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 3 && fields[1] == "wifi" && fields[2] == "connected" {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("активное Wi-Fi соединение не найдено")
}

// getConnectionName получает имя соединения для указанного интерфейса
func getConnectionName(interfaceName string) (string, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,DEVICE", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 2 && fields[1] == interfaceName {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("имя соединения для интерфейса %s не найдено", interfaceName)
}

// setDNS изменяет DNS сервера для указанного соединения
func setDNS(connectionName string, dnsServers []string) error {
	dnsString := strings.Join(dnsServers, ",")
	cmd := exec.Command("nmcli", "connection", "modify", connectionName, "ipv4.dns", dnsString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при выполнении команды: %v\nВывод: %s", err, output)
	}

	// Применяем изменения
	cmd = exec.Command("nmcli", "connection", "up", connectionName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при выполнении команды: %v\nВывод: %s", err, output)
	}

	return nil
}

