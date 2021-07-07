package TinyDatabase

import "encoding/binary"

// HeaderSize :size of KeySize,ValueSize and Tag(uint32,uint32 and uint16)
const HeaderSize = 10

// Status of Tag(Unit)
const (
	PUT uint16 = 0
	DEL	uint16 = 1
)

//Unit :storage unit of TinyDatabase
type Unit struct{
	Key,Value []byte
	KeySize,ValueSize uint32
	Tag uint16
}

// NewUnit :It creates a new Unit.
func NewUnit(Key,Value []byte,Tag uint16)*Unit{
	return &Unit{
		Key:Key,
		Value:Value,
		KeySize: uint32(len(Key)),
		ValueSize: uint32(len(Value)),
		Tag:Tag,
	}
}

// GetSize :It returns the size of a Unit.
func (u *Unit) GetSize() int64{
	return int64(u.KeySize+u.ValueSize+HeaderSize)
}

// Encode :It encodes the Unit to bytes.
func (u *Unit) Encode()([]byte,error){
	buf:=make([]byte,u.GetSize())
	binary.BigEndian.PutUint32(buf[0:4],u.KeySize)
	binary.BigEndian.PutUint32(buf[4:8],u.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10],u.Tag)
	copy(buf[HeaderSize:HeaderSize+u.KeySize],u.Key)
	copy(buf[HeaderSize+u.KeySize:],u.Value)
	return buf,nil
}

//Decode :It decodes the information of KeySize,ValueSize and Tag.
func Decode(buf []byte) (*Unit,error){
	keysize:=binary.BigEndian.Uint32(buf[0:4])
	valuesize:=binary.BigEndian.Uint32(buf[4:8])
	tag:=binary.BigEndian.Uint16(buf[8:10])
	//key:=make([]byte,keysize)
	//value:=make([]byte,valuesize)
	//copy(key,buf[HeaderSize:HeaderSize+keysize])
	//copy(value,buf[HeaderSize+keysize:])
	return &Unit{
		//Key: key,
		//Value: value,
		KeySize: keysize,
		ValueSize: valuesize,
		Tag: tag},
		nil
}