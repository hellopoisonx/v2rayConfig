package yd

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"v2rayConfig/dc"
)

type Server struct {
	Email    string `json:"email"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Level    int    `json:"level"`
}
type Settings struct {
	Servers []Server `json:"servers"`
}
type SocketOptions struct {
	Mark                   int    `json:"mark"`
	TCPFastOpen            bool   `json:"tcpFastOpen"`
	TCPFastOpenQueueLength int    `json:"tcpFastOpenQueueLength"`
	TProxy                 string `json:"tproxy"`
	TCPKeepAliveInterval   int    `json:"tcpKeepAliveInterval"`
}

type StreamSettings struct {
	Network       string        `json:"network"`
	Security      string        `json:"security"`
	TLSSettings   struct{}      `json:"tlsSettings"`
	TCPSettings   struct{}      `json:"tcpSettings"`
	KCPSettings   struct{}      `json:"kcpSettings"`
	WSSettings    struct{}      `json:"wsSettings"`
	HTTPSettings  struct{}      `json:"httpSettings"`
	QUICSettings  struct{}      `json:"quicSettings"`
	DSSettings    struct{}      `json:"dsSettings"`
	GRPCSettings  struct{}      `json:"grpcSettings"`
	SocketOptions SocketOptions `json:"sockopt"`
}
type ProxySettings struct {
	Tag            string `json:"tag"`
	TransportLayer bool   `json:"transportLayer"`
}
type Mux struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency"`
}
type Node struct {
	SendThrough    string         `json:"sendThrough"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	Tag            string         `json:"tag"`
	StreamSettings StreamSettings `json:"streamSettings"`
	ProxySettings  ProxySettings  `json:"proxySettings"`
	Mux            Mux            `json:"mux"`
}

type Outbounds struct {
	Outbounds []Node `json:"outbounds"`
}

func YieldConfig() {
	outbounds := &Outbounds{}
	src := dc.GetSub()
	ss, err := dc.Base64ToSs(string(src))
	if err != nil {
		log.Fatalln("base64解码至ss失败", err)
	}
	for i, v := range ss {
		methodPasswdStart := strings.Index(v, "/")
		methodPasswdEnd := strings.Index(v, "@")
		methodPasswd, err := base64.RawURLEncoding.DecodeString(v[methodPasswdStart+1 : methodPasswdEnd])
		if err != nil {
			log.Fatalln("解码加密方法/密码失败", string(methodPasswd), err)
		}
		methodAndPasswd := strings.Split(string(methodPasswd), ":")
		method := methodAndPasswd[0]
		passwd := methodAndPasswd[1]
		ipStart := methodPasswdEnd
		ipEnd := strings.Index(v, ":")
		ip := v[ipStart+1 : ipEnd]
		portStart := ipEnd
		portEnd := strings.Index(v, "#")
		port, err := strconv.Atoi(v[portStart+1 : portEnd])
		if err != nil {
			log.Fatalln(err)
		}
		outbounds.Outbounds = append(outbounds.Outbounds, Node{})
		outbounds.Outbounds[i].SendThrough = ip
		outbounds.Outbounds[i].Settings.Servers = append(outbounds.Outbounds[i].Settings.Servers, Server{})
		outbounds.Outbounds[i].Settings.Servers[0].Address = ip
		outbounds.Outbounds[i].Settings.Servers[0].Port = port
		outbounds.Outbounds[i].Settings.Servers[0].Method = method
		outbounds.Outbounds[i].Settings.Servers[0].Password = passwd
	}
	encode, err := json.MarshalIndent(outbounds, "", "    ")
	if err != nil {
		log.Fatalln("json转换出现问题", err)
	}
	log.Println(string(encode))
}
