package zkClient

import (
	"testing"
	"time"
)

func Test_CreatePath(t *testing.T) {  
	zkclient,err := New([]string{"192.168.101.161:2181"}, time.Second*10)
	if err !=nil {
		t.Fatalf(err.Error())
		return
	}

	 ex := zkclient.CreatePath("/grs/pay/configs/provider/message", `[{"msg":"msg"}]`)
	if  ex!=nil {
		t.Fatalf(err.Error())
		return
	}
    path,errx := zkclient.CreateSeqNode("/grs/pay/services/pay.get/providers/192", `[{"msg":"msg"}]`)
	if errx!=nil {
		t.Fatalf(errx.Error())
		return
	} 
    
   zkclient.UpdateValue(path,"1234567890")
   if v,_:=zkclient.GetValue(path);v!="1234567890"{
       t.Fatalf("节点修改失败")
   }
   
       
}






