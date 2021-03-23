package security

/*
TS 33.501 6.4.3.1
 COUNT (32 bits) := 0x00 || NAS COUNT (24 bits)
 NAS COUNT (24 bits) := NAS OVERFLOW (16 bits) || NAS SQN (8 bits)
*/
type Count uint32

func (c *Count) maskTo24Bits() {
	*c &= 0x00ffffff
}

func (c *Count) ToUint32() uint32 {
	return uint32(*c)
}

func (c *Count) Set(overflow uint16, sqn uint8) {
	c.SetOverflow(overflow)
	c.SetSQN(sqn)
}

func (c *Count) AddOne() {
	*c++
	c.maskTo24Bits()
}

func (c *Count) SQN() uint8 {
	return uint8(*c & 0x000000ff)
}

func (c *Count) SetSQN(sqn uint8) {
	*c = Count(*c&0x00ffff00) | Count(sqn)
}

func (c *Count) Overflow() uint16 {
	return uint16((*c & 0x00ffff00) >> 8)
}

func (c *Count) SetOverflow(overflow uint16) {
	*c = Count(*c&0x000000ff) | Count(uint32(overflow)<<8)
}
