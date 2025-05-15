package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"log"
	"net"
	"time"
)

// IP IP工具类
var IP = new(ip)

type ip struct{}

// GetExternalIP 获取本机外网IP
func (u ip) GetExternalIP() (string, error) {
	faces, err := net.Interfaces()
	if err != nil {
		return "unknown", err
	}
	for _, iFace := range faces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loop back interface
		}
		address, wr := iFace.Addrs()
		if wr != nil {
			return "unknown", wr
		}
		for _, addr := range address {
			ip := u.getIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "unknown", errors.New("connected to the network")
}

// getIPFromAddr 根据网卡信息获取IP地址
func (u ip) getIPFromAddr(addr net.Addr) net.IP {
	var netIP net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		netIP = v.IP
	case *net.IPAddr:
		netIP = v.IP
	}
	if netIP == nil || netIP.IsLoopback() {
		return nil
	}
	netIP = netIP.To4()
	if netIP == nil {
		return nil // not an ipv4 address
	}

	return netIP
}

// GetOutboundIP get preferred outbound ip of this machine
func (u ip) GetOutboundIP() string {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(8, 8, 8, 8),
		Port: 53,
	})
	if err != nil {
		log.Fatal(err)
	}
	localAddr := socket.LocalAddr().(*net.UDPAddr)
	_ = socket.Close()
	return localAddr.IP.String()
}

// GetLocalIP 获取服务器IP
func (u ip) GetLocalIP() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range address {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return "", errors.New("unknown ip")
}

// IPRegionData 基于IP获取地理位置信息
type IPRegionData struct {
	IP       string `json:"ip"`
	Pro      string `json:"pro" `
	ProCode  int    `json:"proCode" `
	City     string `json:"city" `
	CityCode int    `json:"cityCode"`
	Addr     string `json:"addr"`
	Err      string `json:"err"`
}

// GetRegionDataByIP 根据IP获取城市等相关信息,基于API
func (u ip) GetRegionDataByIP(ip string) (res *IPRegionData) {
	res = &IPRegionData{
		IP:       ip,
		Pro:      "",
		ProCode:  0,
		City:     "",
		CityCode: 0,
		Addr:     "",
		Err:      "",
	}
	url := fmt.Sprintf("http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=%s", ip)
	resp, err := g.Client().Timeout(10*time.Second).Get(context.Background(), url)
	if err != nil {
		return
	}

	defer resp.Close()
	src := resp.ReadAllString()
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	if wr := gconv.Struct(tmp, &res); wr != nil {
		return
	}
	return
}
