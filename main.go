package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"time"
	"github.com/spf13/viper"
	"strings"
)

type Response struct {
	Status bool `json:"status"`
}

func main() {
	viper.SetConfigName("config") // ชื่อ config file
	viper.AddConfigPath(".") // ระบุ path ของ config file
	viper.AutomaticEnv() // อ่าน value จาก ENV variable
	// แปลง _ underscore ใน env เป็น . dot notation ใน viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// อ่าน config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

    ipAddr:=getIp()
	mode:=viper.GetString("app.mode")
	var api string
	if mode=="developement" {
		api = fmt.Sprintf("%s/ip=%s", viper.GetString("app.apitest"),ipAddr)
	}else{
		api = fmt.Sprintf("%s/ip=%s", viper.GetString("app.api"),ipAddr)
	}
	url:=viper.GetString("app.url")
    result,err :=getApi(api)
    if err !=nil{
        os.Exit(1)
    }
    if !result.Status {
        fmt.Println("status is false")
        os.Exit(1)
    }
	openUrl(url)
}
func openUrl(url string){
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
func getIp() string{
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
func getApi(url string) (Response,error) {
    client := http.Client{
        Timeout: 5 * time.Second,
    }
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("No response from request")
		return Response{},err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		fmt.Println("Can not get Response")
		return Response{},err
	}
	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return Response{},err
	}
	fmt.Println(result)
    return result,nil

}
