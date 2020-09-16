package main

import (
	"encoding/binary"
	"fmt"
	"time"

	//"unsafe"
	"os"
)

type Temp struct {
	a1 int8
	a2 int16
	a3 int32
	a4 int64
}

func main(){

	var temp = &Temp{
		100,100,100,100,
	}

	//var byteSlice []byte = *(*[]byte)(unsafe.Pointer(temp))
	//fmt.Println(byteSlice)
	write(temp)
	time.Sleep(3*time.Second)
	read()

}

func write(obj *Temp){
	file , err := os.OpenFile("file.bin", os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	if err := binary.Write(file, binary.BigEndian, obj); err != nil{
		fmt.Println("write in file error: ", err)
	}else{
		fmt.Println("success write to file")
	}
}

func read(){
	var obj = &Temp{}
	file, err := os.OpenFile("file.bin", os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	if err := binary.Read(file, binary.BigEndian, &obj); err != nil {
		fmt.Println("read from file error: ", err)
	}else{
		fmt.Println("success read from file")
	}
	fmt.Printf("%+v\n", *obj)
}
