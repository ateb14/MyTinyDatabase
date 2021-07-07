package TinyDatabase

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type TinyDatabase struct{
	Index map[string]int64
	Ufile *UnitFile
	DirPath string
	rwmutex sync.RWMutex
}

//Open :It opens a TinyDatabase and loads it into RAM.
func Open(dirPath string)(*TinyDatabase,error){
	//If it doesn't exist,just make one.
	if _,err:=os.Stat(dirPath);os.IsNotExist(err){
		if err:=os.MkdirAll(dirPath,os.ModePerm);err!=nil{
			return nil,err
		}
	}

	//Loading the UnitFile from ROM to RAM
	Ufile,err:=NewFile(dirPath)
	if err!=nil{
		log.Fatal(err)
		return nil,err
	}
	db:=&TinyDatabase{
		Index:make(map[string]int64),
		Ufile: Ufile,
		DirPath: dirPath,
	}
	db.IndexInit(Ufile)
	return db,nil
}

// IndexInit Loading the Indexes from UnitFile to RAM
func (db *TinyDatabase) IndexInit(Ufile *UnitFile){
	if Ufile==nil{
		return
	}
	var offset int64 = 0
	for{
		unit,err:=Ufile.Read(offset)
		if err!=nil{
			//Loading finished(EOF)
			if err == io.EOF{
				break
			}
			return
		}
		if unit.Tag==PUT{
			//setting the Index in RAM
			db.Index[string(unit.Key)]=offset
		}
		offset+=unit.GetSize()
	}
	return
}

//Put :It writes the data into the TinyDatabase.
func (db *TinyDatabase) Put(key,value []byte) (err error){
	if len(key)==0{
		return
	}
	db.rwmutex.Lock()
	defer db.rwmutex.Unlock()

	offset:=db.Ufile.Offset
	//packaging the information
	unit:=NewUnit(key,value,PUT)
	//Writing the information to ROM
	err=db.Ufile.Write(unit)
	//Writing the Index to RAM
	db.Index[string(key)]=offset
	fmt.Println("Offset:",offset,"next offset:",db.Ufile.Offset)
	return
}

//Get :It gets the information(value) from ROM.
func (db *TinyDatabase) Get(key []byte) (val []byte,err error) {
	if len(key) == 0 {
		fmt.Println("Empty key!")
		return
	}

	db.rwmutex.RLock()
	defer db.rwmutex.RUnlock()

	//fetching the Index information from RAM
	offset, ok := db.Index[string(key)]
	if !ok {
		fmt.Println("Not found!")
		return
	}

	//fetching the information(value) from ROM
	u, err := db.Ufile.Read(offset)
	if err != nil && err != io.EOF {
		log.Fatal(err)
		return
	}//当前bug：在已有非空文件执行Put后再执行Get后会读到EOF
	if u!=nil {
		val = u.Value
	}else{
		fmt.Println("Empty file!")
	}
	return
}

//Del :It deletes the data in ROM.
func (db *TinyDatabase) Del(key []byte) (err error){
	if len(key)==0 {
		return
	}
	db.rwmutex.Lock()
	defer db.rwmutex.Unlock()

	//fetching the Index information from RAM
	_,ok:=db.Index[string(key)]
	if !ok{
		return
	}

	//
	u:=NewUnit(key,nil,DEL)
	err = db.Ufile.Write(u)
	if err!=nil{
		return
	}

	//deleting the key in hashmap in RAM
	delete(db.Index,string(key))
	return
}

//Update : It clears the redundant units(DEL) in the database.
func (db *TinyDatabase) Update() error{
	//Ignoring the empty database
	if db.Ufile.Offset==0{
		return nil
	}

	var newUnits []*Unit
	var offset int64 = 0

	//Reading the Units in the database and filtering the effective ones
	for{
		unit,err:=db.Ufile.Read(offset)
		if err!=nil{
			if err==io.EOF{
				break
			}
			return err
		}
		//Appending the latest and matched Units to the newUnits slice
		if OFFSET,ok:=db.Index[string(unit.Key)];ok&&OFFSET==offset{
			newUnits=append(newUnits,unit)
		}
		offset+=unit.GetSize()
	}

	if len(newUnits)>0{
		//creating temporary files
		tempUfile,err:=NewTempFile(db.DirPath)
		if err!=nil{
			return err
		}
		defer os.Remove(tempUfile.File.Name())

		//writing the effective Units into the temporary file
		for _,unit:=range newUnits{
			Toffset:=tempUfile.Offset
			err:=tempUfile.Write(unit)
			if err!=nil {
				return err
			}

			//Updating the Index
			db.Index[string(unit.Key)]=Toffset
		}

		//Deleting the old data
		os.Remove(db.Ufile.File.Name())

		//Changing the temporary data to new data
		os.Rename(tempUfile.File.Name(),db.DirPath+string(os.PathSeparator)+Filename)

		//Redirecting
		db.Ufile=tempUfile
	}
	return nil
}