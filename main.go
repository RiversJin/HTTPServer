package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func myLogger() func(r *http.Request) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("Work directory: %v\n", pwd)
	return func(req *http.Request) {
		path, err := url.QueryUnescape(req.URL.Path)
		if err != nil {
			log.Print(err)
		}
		log.Printf("From %v Request for %v", req.RemoteAddr, path)
	}
}
func main() {
	Ips, _ := getClientIp()
	fmt.Println("Please Input the port (Default 8080, 0 means any).")
	var _port int32
	_, err := fmt.Scanf("%d", &_port)
	if err != nil || _port > 0xffff {
		_port = 8080
	}
	Port := strconv.Itoa(int(_port))
	Addrs := make([]string, len(Ips))
	for i := 0; i < len(Ips); i++ {
		ip := Ips[i]
		Addrs[i] = fmt.Sprintf("%s:%s", ip, Port)
	}
	fileServer := FileHandlerWithLog(http.Dir("."), myLogger())
	log.Println("Starting to server, avaliable url:")
	for _, addr := range Addrs {
		log.Printf("http://%s\n", addr)
	}
	err = http.ListenAndServe(":"+Port, fileServer)
	log.Fatal(err)
}
func getClientIp() ([]string, error) {
	broadcastIps := []string{}
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return broadcastIps, err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
			if ipnet.IP.To4() == nil {
				broadcastIps = append(broadcastIps, fmt.Sprintf("[%s]", ipnet.IP.String()))
			} else {
				broadcastIps = append(broadcastIps, ipnet.IP.String())
			}
		}
	}

	if len(broadcastIps) <= 0 {
		return []string{}, errors.New("can not find the client ip address")
	}

	return broadcastIps, nil
}
