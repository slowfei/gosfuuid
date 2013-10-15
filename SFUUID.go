//	Copyright 2013 slowfei And The Contributors All rights reserved.
//	Software Source Code License Agreement (BSD License)

//	Universal Unique IDentifier
//	Reference Implementation https://code.google.com/p/go-uuid/
//
//	email		slowfei@foxmail.com
//	createTime 	2013-8-30
//	updateTime	2013-9-28
//
package SFUUID

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/slowfei/gosfcore/utils/strings"
	"strings"
	"time"
)

var (
	_globalRander = rand.Reader // 全局使用的一个随机种子
)

//	Universal Unique IDentifier is 128 bit (16 byte)
type UUID []byte

//	version 1 new uuid
func NewUUID() UUID {
	return newVersion1()
}

//	version 1 variant node id use public ip
func NewIPUUID() UUID {
	return newIPUUID()
}

//	version 1 custom nodeId byte >= 6
func NewUUIDByNodeID(nodeId []byte) UUID {
	return newVersion1ByNodeID(nodeId)
}

//	version 2 new dce uuid not implement

//	version 3 new md5 uuid
func NewMD5(space UUID, data []byte) UUID {
	if 0 == len(space) || 0 == len(data) {
		return nil
	}
	return newHash(md5.New(), space, data, 3)
}

//	version 4 new random uuid
func NewRandomUUID() UUID {
	return newRandomUUID()
}

//	version 5 new sha1 uuid
func NewSHA1(space UUID, data []byte) UUID {
	if 0 == len(space) || 0 == len(data) {
		return nil
	}
	return newHash(sha1.New(), space, data, 5)
}

//	pares string uuid, error format return nil
//	format	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//	format	urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func Parse(s string) UUID {
	if len(s) == 36+9 {
		if strings.ToLower(s[:9]) != "urn:uuid:" {
			return nil
		}
		s = s[9:]
	} else if len(s) != 36 {
		return nil
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return nil
	}
	uuid := make([]byte, 16)
	for i, x := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		if v, ok := SFStringsUtil.Xtob(s[x:]); !ok {
			return nil
		} else {
			uuid[i] = v
		}
	}
	return uuid
}

func ParseBase64(b64 string) UUID {
	if 22 == len(b64) {

		data, err := base64.URLEncoding.DecodeString(b64 + "==")
		if nil == err {
			return UUID(data)
		}
	}
	return nil
}

func ParseBase64Byte(b64 []byte) UUID {
	if 22 == len(b64) {
		data := make([]byte, 18)
		//	由于在加密的时候除去了结尾的两个(==)符号，所以现在要加上
		_, err := base64.URLEncoding.Decode(data, append(b64, []byte{0x3d, 0x3d}...))
		if nil == err {
			return UUID(data[0:16])
		}
	}
	return nil
}

// Equal returns true if uuid1 and uuid2 are equal.
func Equal(uuid1, uuid2 UUID) bool {
	return bytes.Equal(uuid1, uuid2)
}

//	format	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) String() string {
	if len(u) == 16 {
		b := []byte(u)
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	}
	return ""
}

//	format	urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) URNString() string {
	if len(u) == 16 {
		b := []byte(u)
		return fmt.Sprintf("urn:uuid:%08x-%04x-%04x-%04x-%012x",
			b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	}
	return ""
}

//	base64 format
func (u UUID) Base64() string {
	if len(u) == 16 {
		b := []byte(u)
		//	去除后两位的(==)符号
		return base64.URLEncoding.EncodeToString(b)[0:22]
	}
	return ""
}

//	base64 byte
func (u UUID) Base64Byte() []byte {
	if len(u) == 16 {
		b := []byte(u)
		buf := make([]byte, base64.URLEncoding.EncodedLen(len(b)))
		base64.URLEncoding.Encode(buf, b)
		//	去除后两位的(==)符号
		return buf[0:22]
	}
	return nil
}

//	uuid byte
func (u UUID) Byte() []byte {
	return []byte(u)
}

//	uuid version
func (u UUID) Vresion() int {
	return int(u[6] >> 4)
}

//	uuid variant
//	returns
//	0 = Invalid UUID
//	1 =	RFC4122
//	2 = Reserved, NCS backward compatibility.
//	3 = Reserved, Microsoft Corporation backward compatibility.
//	4 = Reserved for future definition.
func (u UUID) Variant() int {
	if len(u) != 16 {
		return 0 //	Invalid UUID
	}
	switch {
	case (u[8] & 0xc0) == 0x80:
		return 1 //	RFC4122
	case (u[8] & 0xe0) == 0xc0:
		return 2 //	Reserved, NCS backward compatibility.
	case (u[8] & 0xe0) == 0xe0:
		return 3 //	Reserved, Microsoft Corporation backward compatibility.
	default:
		return 4 //	Reserved for future definition.
	}
}

//	version 1 and version 2 use
//	return the create time
//	format error return false
func (u UUID) Time() (time.Time, bool) {
	if 1 == u.Vresion() || 2 == u.Vresion() {
		if len(u) == 16 {
			t := int64(binary.BigEndian.Uint32(u[0:4]))
			t |= int64(binary.BigEndian.Uint16(u[4:6])) << 32
			t |= int64(binary.BigEndian.Uint16(u[6:8])&0xfff) << 48

			sec := int64(t)
			nsec := (sec % 10000000) * 100
			sec /= 10000000
			return time.Unix(sec, nsec), true
		}
	}
	return time.Time{}, false
}

//	version 1 and version 2 use
//	return the uuid nodeId
//	format error return false
func (u UUID) NodeID() ([]byte, bool) {
	if 1 == u.Vresion() || 2 == u.Vresion() {
		if len(u) == 16 {
			node := make([]byte, 6)
			copy(node, u[10:])
			return node, true
		}

	}
	return nil, false
}

//	version 1 and version 2 use
//	return the uuid clock sequence
//	format error return false
func (u UUID) ClockSequence() (int, bool) {
	if 1 == u.Vresion() || 2 == u.Vresion() {
		if len(u) == 16 {
			return int(binary.BigEndian.Uint16(u[8:10])) & 0x3fff, true
		}
	}
	return 0, false

}

//	uuid error info
type UUIDError struct {
	Message string
}

func NewUUIDError(format string, args ...interface{}) *UUIDError {
	return &UUIDError{fmt.Sprintf(format, args...)}
}

func (err *UUIDError) Error() string {
	return err.Message
}
