The UUID Implementation
========

Universal Unique IDentifier<br/>
Reference Implementation https://code.google.com/p/go-uuid/<br/>

gosfuuid一些代码是参考了go-uuid进行实现的，增加了base64的uuid转换字符串，达到22位的字符。并将代码进行了一些优化，调用更加的方便。

#### Install And Use

	go get github.com/slowfei/gosfuuid

```golang
impurt uuid "github.com/slowfei/gosfuuid"
impurt "fmt"

func main(){
	uid = uuid.NewRandomUUID()
	base64Str := uid.Base64()
	fmt.Println(base64Str)
}

```

##
#### 使用协议 [LICENSE](https://github.com/slowfei/gosfuuid/blob/master/LICENSE)

Software Source Code License Agreement (BSD License)

###
