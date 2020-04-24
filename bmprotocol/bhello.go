package bmprotocol

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/tools"
)

type BMHello struct {

}

func NewBMHello() *BMHello {
	return &BMHello{}
}

func (bmh *BMHello)Pack() ([]byte, error) {
	bmact:= translayer.NewBMTL(translayer.HELLO,nil)

	return bmact.Pack()
}

func (bmh *BMHello)UnPack(data []byte) (int,error) {
	//nothing todo
	return 0,nil
}


type BMHelloACK struct {
	sn []byte
}

func NewBMHelloACK(sn []byte) *BMHelloACK {
	return &BMHelloACK{}
}



func (bmha *BMHelloACK)Pack() ([]byte,error)  {
	var barr []byte

	bufl:=translayer.UInt16ToBuf(uint16(len(bmha.sn)))

	barr = append(barr,bufl...)
	barr = append(barr,bmha.sn...)

	bmact:= translayer.NewBMTL(translayer.HELLO_ACK,barr)

	return bmact.Pack()
}

func (bmha *BMHelloACK)UnPack(data []byte) (int,error) {

	if len(data)<2{
		return 0,errors.New("Not a HELLO ACK data")
	}

	l:=binary.BigEndian.Uint16(data)

	bmha.sn = data[2:]

	if l != uint16(len(bmha.sn)){
		return 0,errors.New("Serial Nunber Error")
	}

	return 2+len(bmha.sn),nil
}


/*
*

pack:
1.NewSendSignature
2.ForSigBuf
3.Calculate signature
4.SetSig
5.Pack
6.Send

unpack:
bmtl:=&BMTransLayer{}
bmtl.UnPack(buf)
if bmtl.typ == SEND_SIGNATURE{
	ss:=&SendSignature{}
	ss.UnPack(buf[bmtl.GetData()])
}

*/

type SendSignature struct {
	sn []byte
	localMailAddr string
	currentTime int64      //Millisecond
	sig []byte
}

func NewSendSignature(sn []byte,localMailAddr string) *SendSignature {
	ss := &SendSignature{}
	ss.sn = sn
	ss.localMailAddr = localMailAddr
	ss.currentTime = tools.GetNowMsTime()

	return ss
}



func (ss *SendSignature)ForSigBuf() []byte  {
	var r []byte

	bufl:=translayer.UInt16ToBuf(uint16(len(ss.sn)))

	r = append(r,bufl...)

	r = append(r,ss.sn...)

	bufl = translayer.UInt16ToBuf(uint16(len(ss.localMailAddr)))

	r = append(r,bufl...)

	r = append(r,[]byte(ss.localMailAddr)...)

	bufl = translayer.UInt64ToBuf(uint64(ss.currentTime))

	r = append(r,bufl...)

	return r
}

func (ss *SendSignature)SetSig(sig []byte)  {
	ss.sig = sig
}

func (ss *SendSignature)Pack() ([]byte,error)  {
	r:=ss.ForSigBuf()

	bufl:=translayer.UInt16ToBuf(uint16(len(ss.sig)))

	r = append(r,bufl...)

	r = append(r,ss.sig...)

	bmact:= translayer.NewBMTL(translayer.SEND_SIGNATURE,r)

	return bmact.Pack()
}

func (ss *SendSignature)UnPack(buf []byte)  (int,error) {
	offset:=0

	if len(buf[offset:])<2{
		return 0,errors.New("Not a SendSignature data")
	}

	lsn:=binary.BigEndian.Uint16(buf[offset:])
	offset += 2
	if len(buf)<offset+int(lsn){
		return 0,errors.New("Not a SendSignature data")
	}
	ss.sn = buf[offset:offset+int(lsn)]
	offset += int(lsn)

	if len(buf[offset:])<2{
		return 0,errors.New("Not a SendSignature data")
	}
	laddr:=binary.BigEndian.Uint16(buf[offset:])
	offset += 2
	if len(buf)<offset+int(laddr){
		return 0,errors.New("Not a SendSignature data")
	}
	ss.localMailAddr = string(buf[offset:offset+int(laddr)])
	offset += int(laddr)
	if len(buf[offset:])<8{
		return 0,errors.New("Not a SendSignature data")
	}
	ss.currentTime = int64(binary.BigEndian.Uint64(buf[offset:]))
	offset += 8

	if len(buf[offset:])<2{
		return 0,errors.New("Not a SendSignature data")
	}

	lsig:=binary.BigEndian.Uint16(buf[offset:])
	offset += 2
	if len(buf)<offset + int(lsig){
		return 0,errors.New("Not a SendSignature data")
	}

	ss.sig = buf[offset:]

	return offset + int(lsig),nil

}

type ValidateSignature struct {
	sn []byte
}


func NewValidSign(sn []byte) *ValidateSignature {
	return &ValidateSignature{}
}

func (vs *ValidateSignature)Pack() ([]byte,error)  {
	var barr []byte

	bufl:=translayer.UInt16ToBuf(uint16(len(vs.sn)))

	barr = append(barr,bufl...)
	barr = append(barr,vs.sn...)

	bmact:= translayer.NewBMTL(translayer.VALIDATE_SIGNATURE,barr)

	return bmact.Pack()
}

func (vs *ValidateSignature)UnPack(data []byte) (int,error) {

	if len(data)<2{
		return 0,errors.New("Not a Validate Signature data")
	}

	l:=binary.BigEndian.Uint16(data)

	vs.sn = data[2:]

	if l != uint16(len(vs.sn)){
		return 0,errors.New("Serial Nunber Error")
	}

	return 2+len(vs.sn),nil
}
