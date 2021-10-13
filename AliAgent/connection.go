package AliAgent

import (
	"net/http"
)


func FetchRemotes() {
	http.Get(`https://staragent-configservice.aliyuncs.com/api/configservice?action=findChannelListForAgent&agentIpList=172.21.129.154&needAllChannels=true&serviceTag=a6b54ef8-d51e-4d88-8d7c-4d421923d2d6&version=2`)
}