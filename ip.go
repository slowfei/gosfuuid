//	Copyright 2013 slowfei And The Contributors All rights reserved.
//	Software Source Code License Agreement (BSD License)

//	UUID Version 1 Variant 通过计算当前时间戳、随机数和本机IP(外网)地址得到uuid
//	考虑到无法链接物联网时获取不了外网的问题，则使panic()爆出异常。
//	由version 1进行变种
//
//	email		slowfei@foxmail.com
//	createTime 	2013-8-31
//	updateTime	2013-9-28
//
package SFUUID

import (
	"encoding/binary"
	"github.com/slowfei/gosfcore/utils/rand"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	//	默认获取网络IP数据的URL
	UUID_URL_IP_API = "http://api.externalip.net/ip/"
)

var (
	_ipNodeId []byte //	ip nodeid
)

func newIPUUID() UUID {
	if nil == _ipNodeId {
		//	这里有可能再多并发的时候会全部涌进，所以最好再程序开始运行时进行设置，或手动设置IP
		SetNetwordIP("")
	}

	//	时间戳
	now := getTimestamp()

	uuid := make([]byte, 16)

	// 15 – 12: TimeLow 时间值的低位
	// 11 – 10: TimeMid 时间值的中位
	// 09 – 08: VersionAndTimeHigh 4位版本号和时间值的高位
	// 07: VariantAndClockSeqHigh 2位变体（ITU-T）和时钟序列高位
	// 06: ClockSeqLow 时钟序列低位
	// 05 – 00: Node 结点

	time_low := uint32(now & 0xffffffff)
	time_mid := uint16((now >> 32) & 0xffff)
	time_hi := uint16((now >> 48) & 0x0fff)
	time_hi |= 0x1000 // Version 1

	binary.BigEndian.PutUint32(uuid[0:], time_low)
	binary.BigEndian.PutUint16(uuid[4:], time_mid)
	binary.BigEndian.PutUint16(uuid[6:], time_hi)
	binary.BigEndian.PutUint16(uuid[8:], _clock_seq)
	uuid[8] = (uuid[8] & 0x3f) | 0x10 // clear and set to set to variant is 10

	copy(uuid[10:], _ipNodeId)

	return uuid
}

//	设置自定义的IP
//	@ip set net.IP
func SetIPNodeId(ip net.IP) {
	if nil != ip {
		_ipNodeId = []byte(ip)
	}
}

//	设置网络IP作为nodeid之用
//	@urlIPApi	获取ip的地址链接，地址访问是可以直接读取ip数据的，不需要任何解析.
//				传递("")则使用默认 UUID_URL_IP_API
//
func SetNetwordIP(urlIPApi string) {

	if nil == _ipNodeId {
		_ipNodeId = make([]byte, 6)
	}

	SFRandUtil.RandBits(_ipNodeId)

	var url string
	if "" != urlIPApi {
		url = urlIPApi
	} else {
		url = UUID_URL_IP_API
	}

	isPanic := true
	res, err := http.Get(url)
	if nil == err {
		defer res.Body.Close()

		if data, err := ioutil.ReadAll(res.Body); err == nil {
			ip := net.ParseIP(string(data))
			if nil != ip && nil != ip.To4() {
				copy(_ipNodeId, []byte(ip.To4()))
				isPanic = false
			}
		}
	}
	if isPanic {
		panic(NewUUIDError("internet link failure, not obtain IP information"))
	}

}
