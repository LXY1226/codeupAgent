package AliAgent

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const defaultName = `staragentd`
const defaultSNName = `staragent_sn`

var remoteConn *net.TCPConn

func InitAgent() {
	type Header struct {
		Hostname           string `json:"hostname"`
		IP                 string `json:"ip"`
		Os                 string `json:"os"`
		OsBit              string `json:"osBit"`
		SN                 string `json:"sn"`
		StarAgentStartTime string `json:"staragentStartTime"`
	}
	type ppfReport struct {
		Body   jsoniter.RawMessage `json:"body"`
		Header Header              `json:"header"`
	}
	var report ppfReport

	{
		report.Body = jsoniter.RawMessage(`[]`)
		sn, err := os.ReadFile(defaultPath)
		if err != nil {

			sn, err = os.ReadFile(SNPath())
			if err != nil {
				fmt.Println("SN不存在，初次使用请粘贴添加命令到此处然后回车")
				var bash, curl, install, SN, Agent, Verify, Region, TimeStamp, Proxy string
				n, err := fmt.Scan(&bash, &curl, &install, &SN, &Agent, &Verify, &Region, &TimeStamp, &Proxy)
				if err != nil {
					log.Panicln("输入有误", err)
				}
				if n != 9 {
					log.Panicln("输入有误", err)
				}
				resp, err := http.Get(SN[1 : len(SN)-2])
				if err != nil {
					log.Panicln("不能获取SN", err)
				}
				sn, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Panicln("不能读取SN", err)
				}
				if string(sn) == "-1" {
					log.Fatalln("签名无效")
				}
				if string(sn) == "-2" {
					log.Fatalln("安装命令已过期，请到机器管理页面重新生成安装命令。")
				}
				err = os.WriteFile(SNPath(), sn, 0644)
				if err != nil {
					log.Fatalln("无法写入SN文件", err)
				}
			}
		}
		report.Header.Hostname, err = os.Hostname()
		if err != nil {
			report.Header.Hostname = "(unknown)"
		}
		report.Header.SN = string(sn)
	}

	{
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if addr.(*net.IPNet).IP.To4() != nil {
					report.Header.IP = addr.String()
					break
				}
			}
		}
	}

	report.Header.Os = runtime.GOOS
	//report.Header.OsBit = runtime.GOARCH
	report.Header.OsBit = "32"

	{
		buf := make([]byte, 512)
		buf = append(buf, `https://staragent-configservice.aliyuncs.com/api/configservice?action=findChannelListForAgent&needAllChannels=true&serviceTag=`...)
		buf = append(buf, report.Header.SN...)
		buf = append(buf, `&agentIpList=`...)
		buf = append(buf, report.Header.IP...)
		for {
			resp, err := http.Get(string(buf))
			if err != nil {
				log.Println("无法获取服务器列表", err, "重试...")
				time.Sleep(2 * time.Second)
				continue
			}
			type ChannelIPPort struct {
				IP   string `json:"ip"`
				Port int    `json:"port"`
			}
			body, err := io.ReadAll(resp.Body)
			var remotes []ChannelIPPort
			jsoniter.Get(body, "result", "channelIPPort").ToVal(&remotes)
			ch := make(chan error)
			wg := sync.WaitGroup{}
			wg.Add(len(remotes))
			for _, remote := range remotes {
				remote := remote
				go func() {
					defer wg.Done()
					conn, err := net.DialTCP("tcp",nil, &net.TCPAddr{IP: net.ParseIP(remote.IP), Port: remote.Port})
					if err != nil {
						return
					}
					if remoteConn == nil {
						atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&remoteConn)), unsafe.Pointer(conn))
						log.Println("已连接至", remote.IP)
					}
				}()
			}
			go func() {
				wg.Wait()
				if remoteConn == nil {
					ch <- errors.New("no address usable")
				}
			}()
			if err = <- ch; err != nil {
				log.Println(err)
			}
			break
		}
	}

	type channelReg struct {
		AgentRegisterStartTime string `json:"agentRegisterStartTime"`
		AgentStartTime         string `json:"agentStartTime"`
		ClientOs               string `json:"clientOs"`
		ClientOsBit            string `json:"clientOsBit"`
		ClientVersion          string `json:"clientVersion"`
		HostName               string `json:"hostName"`
		IP                     string `json:"ip"`
		IPList                 string `json:"ipList"`
		Role                   string `json:"role"`
		ServerStartTime        string `json:"serverStartTime"`
		ServiceTag             string `json:"serviceTag"`
		Status                 string `json:"status"`
		UUID                   string `json:"uuid"`
	}
}

func sendPacket(msgType,trafficType string, body string) {

}