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

type fileHandlerWithLog struct {
	handler http.Handler
	logger  func(*http.Request)
}

func FileHandlerWithLog(root http.FileSystem, logger func(*http.Request)) http.Handler {
	return &fileHandlerWithLog{http.FileServer(root), logger}
}

func (fhw *fileHandlerWithLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fhw.logger(r)
	fhw.handler.ServeHTTP(w, r)
}
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
	Ip, _ := getClientIp()
	fmt.Println("Please Input the port (Default 8080, 0 means any).")
	var _port int32
	_, err := fmt.Scanf("%d", &_port)
	if err != nil || _port > 0xffff {
		_port = 8080
	}
	Port := strconv.Itoa(int(_port))
	Addr := fmt.Sprintf("%s:%s", Ip, Port)
	fileServer := FileHandlerWithLog(http.Dir("."), myLogger())
	log.Printf("Starting at http://%s", Addr)
	err = http.ListenAndServe(":"+Port, fileServer)
	log.Fatal(err)
}
func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("can not find the client ip address")
}