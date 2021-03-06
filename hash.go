//	Copyright 2013-2015 slowfei And The Contributors All rights reserved.
//
//	Software Source Code License Agreement (BSD License)
//
//  Create on 2013-09-02
//  Update on 2015-08-14
//  Email  slowfei@nnyxing.com
//  Home   http://www.slowfei.com

/***UUID Version 3\5
　　UUID Version 3(md5) and Version 5(sha1)
<br/>
Reference Implementation https://code.google.com/p/go-uuid/
*/

// UUID Version 3\5
package SFUUID

import (
	"hash"
)

// Well known Name Space IDs and UUIDs
// http://www.ietf.org/rfc/rfc4122.txt
var (
	NameSpace_DNS  = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NameSpace_URL  = UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NameSpace_OID  = UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NameSpace_X500 = UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
)

//	将需要的数据转化成UUID识别码，可以使用md5或sha1进行hash的编码转换
//	在同样的命名空间和相同的数据UUID是一样的。
//	如果不同的命名空间和相同的数据UUID是不会一样的
//	默认使用的是rfc4122定义的命名空间
func newHash(h hash.Hash, space UUID, data []byte, version int) UUID {
	h.Reset()
	h.Write(space)
	h.Write([]byte(data))
	s := h.Sum(nil)
	uuid := make([]byte, 16)
	copy(uuid, s)
	uuid[6] = (uuid[6] & 0x0f) | uint8((version&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // RFC 4122 variant
	return uuid
}
