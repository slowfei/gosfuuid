//	Copyright 2013 slowfei And The Contributors All rights reserved.
//	Software Source Code License Agreement (BSD License)

//	UUID Version 4 random uuid
//	Reference Implementation https://code.google.com/p/go-uuid/
//
//	email		slowfei@foxmail.com
//	createTime 	2013-8-30
//	updateTime	2013-9-26
//
package SFUUID

import (
	"github.com/slowfei/gosfcore/utils/rand"
)

//	创建一个随机数的uuid
func newRandomUUID() UUID {
	//	随机数的UUID没有什么特别，就是随机数

	rb := make([]byte, 16)
	SFRandUtil.RandBits(rb)

	rb[6] = (rb[6] & 0x0f) | 0x40 // clear and set to Version 4
	rb[8] = (rb[8] & 0x3f) | 0x80 // clear and set to set to variant is 10
	return UUID(rb)
}
