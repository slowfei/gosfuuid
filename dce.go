//	Copyright 2013-2015 slowfei And The Contributors All rights reserved.
//
//	Software Source Code License Agreement (BSD License)
//
//  Create on 2013-08-31
//  Update on 2015-08-14
//  Email  slowfei@nnyxing.com
//  Home   http://www.slowfei.com

/***UUID Version 2 DCE 1.1
　　（Distributed Computing Environment）uuid
<br/>
Reference Implementation https://code.google.com/p/go-uuid/
*/

// UUID Version 2 DCE 1.1
package SFUUID

import (
	"encoding/binary"
	"os"
)

func newDCESecurity() UUID {
	//	TODO 生成出来的数据感觉有点不对头，os.Getuid()获取的是基本的一个固定值，会把时间戳的前4位置换为POSIX的UID或GID，
	//	这样产生在同一台机的UUID几乎都是一样的，这个版本的 uuid不知道用于什么方向，暂时不实现先了
	uuid := newVersion1()
	if uuid != nil {
		uuid[6] = (uuid[6] & 0x0f) | 0x20 // Version 2
		uuid[9] = byte(0)
		binary.BigEndian.PutUint32(uuid[0:], uint32(os.Getpid()))
	}
	return uuid
}
