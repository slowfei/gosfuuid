//	Copyright 2013-2015 slowfei And The Contributors All rights reserved.
//
//	Software Source Code License Agreement (BSD License)
//
//  Create on 2013-08-31
//  Update on 2015-08-14
//  Email  slowfei@nnyxing.com
//  Home   http://www.slowfei.com
//

/***UUID Version 1
　　通过计算当前时间戳、随机数和机器MAC地址得到uuid
<br/>
Reference Implementation https://code.google.com/p/go-uuid/
*/

//	UUID Version 1
package SFUUID

import (
	"encoding/binary"
	"github.com/slowfei/gosfcore/utils/rand"
	"net"
	"sync"
	"time"
)

var (
	_interfaces    []net.Interface // net.Interfaces()  system's network interfaces.
	_interfaceName string          // network interface name
	_nodeId        []byte          // version 1 hardware address
	_rwmutex       sync.RWMutex    //
	_lasttime      uint64          // last time we returned
	_clock_seq     uint16          // clock sequence for this run
)

func init() {
	//	为避免一开始线程涌进，所以先设置基本参数
	// fmt.Println("version init")
	//	设置时间序
	SetClockSequence(-1)
	//	获取网络接口的信息，设置成nodeId
	SetNodeInterface("")
}

//  new uuid version 1
//	uuid version 1 是基于计算当前时间戳、随机数和mac地址得到的。
//	version 1可以保证在全球范围的唯一性，但与此同时，使用MAC地址会带来安全性问题。
//	如果担心安全问题可以使用version 2 dce的版本，version 2是基于version 1产生的.
func newVersion1() UUID {
	return newVersion1ByNodeID(nil)
}

//	@nodeId	自定义网络节点标识
func newVersion1ByNodeID(nodeId []byte) UUID {
	//	时间戳
	now, clockSeq := getTimestamp()

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
	binary.BigEndian.PutUint16(uuid[8:], clockSeq)
	if nil != nodeId {
		if len(nodeId) >= 6 {
			copy(uuid[10:], nodeId)
		} else {
			panic(NewUUIDError("use defined uuid node id error:%v", nodeId))
		}
	} else {
		copy(uuid[10:], _nodeId)
	}

	return uuid
}

//	获取时间戳
func getTimestamp() (int64, uint16) {
	//	连续生成两个UUID的时间至少要间隔100ns. 考虑到并发的问题可能会出现相同的时间戳，使用锁来控制。
	_rwmutex.RLock()
	defer _rwmutex.RUnlock()

	t := time.Now()

	//	由于uuid 16 byte存储时间的byte有限，所以需要除去100ns
	//	1378145512981336289 = 13781455129813362
	//	后续还原时间的时候会有100ns的误差
	now := uint64(t.UnixNano() / 100)

	//	对基于时间的UUID版本，时间序列用于避免因时间向后设置或节点值改变可能造成的UUID重复，
	//	对基于名称或随机数的版本同样有用：目的都是为了防止UUID重复。
	if now <= _lasttime {
		_clock_seq = ((_clock_seq + 1) & 0x3fff) | 0x8000
	}
	_lasttime = now

	return int64(now), _clock_seq
}

func ClockSequence() int {
	if _clock_seq == 0 {
		SetClockSequence(-1)
	}
	return int(_clock_seq & 0x3fff)
}

//	设置时间序，设置低于 14 bit
//	@seq 	如果为-1将使用随机数产生一个新的时间序
func SetClockSequence(seq int) {
	if -1 == seq {
		// random clock sequence

		b := make([]byte, 2)
		SFRandUtil.RandBits(b)

		seq = int(b[0])<<8 | int(b[1])
	}
	old_seq := _clock_seq
	_clock_seq = uint16(seq&0x3fff) | 0x8000 // set our variant
	if old_seq != _clock_seq {
		_lasttime = 0
	}
}

//	设置网络接口nodeId，根据自己的需求传递interface name进行设置
//	使用于 uuid version 1
//
//	@name	传递(name="")如果获取补了网络接口信息会使用随机数代替
//	@return bool  返回是否设置成功
//	@return error 错误信息
func SetNodeInterface(name string) (bool, error) {

	var err error
	if nil == _interfaces {
		_interfaces, err = net.Interfaces()
		if nil != err && "" == name {
			return false, err
		}
	}

	for _, ifs := range _interfaces {
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) {
			if setNodeID(ifs.HardwareAddr) {
				_interfaceName = ifs.Name
				return true, nil
			}
		}
	}

	//	考虑到网络或其他因素，可能获取不了mac信息，使用随机数做为nodeId
	//	随机数设置只在name等于空的情况下进行
	if "" == name {

		if nil == _nodeId {
			_nodeId = make([]byte, 6)
		}
		SFRandUtil.RandBits(_nodeId)

		return true, nil
	}

	return false, nil
}

//	返回当前uuid version 1所使用的网络接口名称
func NodeInterface() string {
	return _interfaceName
}

//	获取当前uuid version 1所使用的nodeid
func NodeID() []byte {
	if _nodeId == nil {
		SetNodeInterface("")
	}
	result := make([]byte, 6)
	copy(result, _nodeId)
	return result
}

//	可以自定义设置uuid的nodeid，id必须大于或等于6 bytes，不过也只获取前 6 bytes
func SetNodeID(id []byte) bool {
	if setNodeID(id) {
		_interfaceName = "UserDefined"
		return true
	}
	return false
}
func setNodeID(id []byte) bool {
	if len(id) < 6 {
		return false
	}
	if _nodeId == nil {
		_nodeId = make([]byte, 6)
	}
	copy(_nodeId, id)
	return true
}
