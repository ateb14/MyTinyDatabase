package TinyDatabase

import (
	"log"
	"os"
)

const Filename = "tinydb.data"
const TempFilename = "tinydb.data.temp"

type UnitFile struct{
	File *os.File
	Offset int64
}

func newInternal (filename string) (*UnitFile,error){
	file,err:=os.OpenFile(filename,os.O_CREATE|os.O_RDWR,0644)
	if err!=nil{
		return nil,err
	}
	stat,err:=os.Stat(filename)
	if err!=nil{
		return nil,err
	}
	return &UnitFile{file,stat.Size()},nil
}
//NewFile :It opens or creates a file in ROM and returns it back
func NewFile (DirPath string) (*UnitFile,error){
	filename:=DirPath+string(os.PathSeparator)+Filename//filename is also the filepath
	return newInternal(filename)
}

//NewTempFile :It creates a temporary file in ROM and returns it back
func NewTempFile(DirPath string) (*UnitFile,error){
	filename:=DirPath+string(os.PathSeparator)+TempFilename
	return newInternal(filename)
}
//Read :It reads the UnitFile from the starting byte offset.
func (uf *UnitFile) Read(offset int64)(u *Unit,err error){
	buf:=make([]byte,HeaderSize)
	_,err=uf.File.ReadAt(buf,offset)
	if err!=nil{
		return
	}
	u,err=Decode(buf)
	if err!=nil{
		return
	}

	//Reading the key
	offset+=HeaderSize
	if u.KeySize>0{
		key:=make([]byte,u.KeySize)
		if _,err=uf.File.ReadAt(key,offset);err!=nil{
			return
		}
		u.Key=key
	}

	//Reading the value
	offset+=int64(u.KeySize)
	if u.ValueSize>0{
		value:=make([]byte,u.ValueSize)
		if _,err=uf.File.ReadAt(value,offset);err!=nil{
			return
		}
		u.Value=value
	}
	return
}

//Write :It writes the Unit into the UnitFile
func (uf *UnitFile) Write(u *Unit)(err error){
	message,err:=u.Encode()
	if err!=nil{
		return err
	}
	_,err=uf.File.Write(message)
	uf.Offset+=u.GetSize()
	if err!=nil{
		log.Fatal(err)
	}
	return
}
