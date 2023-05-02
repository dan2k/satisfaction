package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Response struct {
	Status bool `json:"status"`
}

func main() {

	viper.SetConfigName("config") // ชื่อ config file
	viper.AddConfigPath(".")      // ระบุ path ของ config file
	viper.AutomaticEnv()          // อ่าน value จาก ENV variable
	// แปลง _ underscore ใน env เป็น . dot notation ใน viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// อ่าน config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	// currentTime := time.Now()
	// c:=fmt.Sprintf("Current Time in String: %s", currentTime.String())
	f, err := os.OpenFile("start.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)
	logger.Println(" start finish")
	ipAddr := getIp()
	mac, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// fmt.Println(mac)
	mode := viper.GetString("app.mode")
	var api string
	if mode == "developement" {
		api = fmt.Sprintf("%s?ip=%s&mac=%s", viper.GetString("app.apitest"), ipAddr, mac)
	} else {
		api = fmt.Sprintf("%s?ip=%s&mac=%s", viper.GetString("app.api"), ipAddr, mac)
	}
	fmt.Println(api)
	url := viper.GetString("app.url")
	result, err := getApi(api)
	if err != nil {
		os.Exit(1)
	}
	if !result.Status {
		fmt.Println("status is false")
		os.Exit(1)
	}
	openUrl(fmt.Sprintf(""+" %s&mac=%s",url,mac))
}
func openUrl(url string) {
	// fmt.Println(url)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "chrome", url)
	case "darwin":
		cmd = exec.Command("open", "-a", "Google Chrome", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error opening URL:", err)
		os.Exit(1)
	}
}
func getIp() string {
	// Get request
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	var ipAddr string

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			// fmt.Println("IPv4: ", ipv4)
			ipAddr = ipv4.String()
			break
		}
	}
	return ipAddr
}
func getApi(url string) (Response, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("No response from request")
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		fmt.Println("Can not get Response")
		return Response{}, err
	}
	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return Response{}, err
	}
	fmt.Println(result)
	return result, nil

}

func getMacAddr() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var as string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = a
			break
		}
	}
	return as, nil
}
