//Just for testing
package main

import (
	"./TinyDatabase"
	"fmt"
	"log"
)

func main(){
	db,err:=TinyDatabase.Open("/Users/wangyiyang/Documents/GoProject/DataBase/database")
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("DirPath is :",db.DirPath,"\noffset is",db.Ufile.Offset)
	key1:=[]byte{65}
	key2:=[]byte{65,66}
	key3:=[]byte{65,67}
	key4:=[]byte{65,66,67,73}
	key5:=[]byte{'t','e','s','t'}
	value1:=[]byte{'o','n','e'}
	value2:=[]byte{'t','w','o'}
	value3:=[]byte{'t','h','r','e','e'}
	value4:=[]byte{'f','o','u','r'}
	//value5:=[]byte{'c','e','n'}
	db.Put(key1,value1)
	db.Put(key2,value2)
	db.Put(key3,value3)
	db.Put(key4,value4)
	fmt.Println(db.Index)
	fmt.Println("\noffset is",db.Ufile.Offset)
	//db.Put(key5,value5)
	//db.Update()
	ret,_:=db.Get(key1)
	fmt.Println(string(ret))
	ret,_=db.Get(key2)
	fmt.Println(string(ret))
	ret,_=db.Get(key3)
	fmt.Println(string(ret))
	ret,_=db.Get(key4)
	fmt.Println(string(ret))
	ret,_=db.Get(key5)
	fmt.Println(string(ret))
	fmt.Println(db.Index)
	fmt.Println("\noffset is",db.Ufile.Offset)
}
