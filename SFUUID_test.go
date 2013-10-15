package SFUUID

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestRandomUUID(t *testing.T) {

	for i := 0; i < 10; i++ {
		uuid := NewRandomUUID()
		fmt.Println(uuid.String())
		fmt.Println(uuid.Base64())
		fmt.Println(uuid.Base64Byte())
		fmt.Println("version:", uuid.Vresion())
		fmt.Println("Variant:", uuid.Variant())
	}

}

func TestMD5UUID(t *testing.T) {
	data := []byte("slowfei")
	space := NameSpace_X500
	uuid := NewMD5(space, data)
	fmt.Println(uuid.String())
	fmt.Println("version:", uuid.Vresion())
	fmt.Println("Variant:", uuid.Variant())
	if uuid.String() != "2e6d58ae-1432-3541-b8b2-d456c76fb43b" {
		t.Fail()
	}
}

func TestSHA1UUID(t *testing.T) {
	data := []byte("slowfei")
	space := NameSpace_X500
	uuid := NewSHA1(space, data)
	fmt.Println(uuid.String())
	fmt.Println("version:", uuid.Vresion())
	fmt.Println("Variant:", uuid.Variant())
	if uuid.String() != "64809754-fd5a-51a2-b0a4-7ad622549ab5" {
		t.Fail()
	}
}

func TestVersion1UUID(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	testU := func(tag string, c chan bool) {
		for i := 0; i < 5; i++ {
			uuid1 := NewUUID()
			fmt.Println("version:", tag, "- ", uuid1.String())
			t, _ := uuid1.Time()
			fmt.Println("version-time:", t)
			nodeId, _ := uuid1.NodeID()
			fmt.Println("version-node:", nodeId)
			fmt.Println("")

		}
		c <- true
	}

	a := make(chan bool)
	b := make(chan bool)
	c := make(chan bool)
	d := make(chan bool)
	e := make(chan bool)
	f := make(chan bool)
	go testU("go1", a)
	go testU("go2", b)
	go testU("go3", c)
	go testU("go4", d)
	go testU("go5", e)
	go testU("go6", f)

	<-a
	<-b
	<-c
	<-d
	<-e
	<-f
	fmt.Println("current node:", NodeID())
}

func TestDCEUUID(t *testing.T) {
	fmt.Println(newDCESecurity().String())
	fmt.Println(newDCESecurity().String())
	fmt.Println(newDCESecurity().String())

	uuid := make([]byte, 20)
	binary.BigEndian.PutUint32(uuid[0:], uint32(os.Getuid()))
	fmt.Println(uuid)
}

func TestParse(t *testing.T) {
	uuidStr := "975081b7-f6f7-1030-af40-20c9d0442301"

	uuid := Parse(uuidStr)

	time, _ := uuid.Time()
	if time.String() != "2013-09-04 02:50:13.8102199 +0800 CST" {
		t.Errorf("uuid time error:%v", time.String())
	}

	nodeId, _ := uuid.NodeID()
	if !bytes.Equal(nodeId, []byte{0x20, 0xc9, 0xd0, 0x44, 0x23, 0x01}) {
		t.Errorf("uuid nodeId error:%v", nodeId)
	}

}

func TestParseBase64(t *testing.T) {
	base64Str := "WmlexAMhT5-si3vNCn3MNQ"
	uuidStr := "5a695ec4-0321-4f9f-ac8b-7bcd0a7dcc35"

	uuid := ParseBase64(base64Str)

	if uuid.String() != uuidStr {
		t.Errorf("uuid base64 parse error:%v", uuid.String())
	}
}

func TestIPUUID(t *testing.T) {

	for i := 0; i < 10; i++ {
		uuid := NewIPUUID()
		fmt.Println(uuid.String())
		fmt.Println(uuid.Variant())
		fmt.Println(uuid.Time())
	}
}
