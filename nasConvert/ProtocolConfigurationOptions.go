package nasConvert

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/free5gc/nas/nasMessage"
)

type ProtocolOrContainerUnit struct {
	ProtocolOrContainerID uint16
	LengthOfContents      uint8
	Contents              []byte
}

type ProtocolConfigurationOptions struct {
	ProtocolOrContainerList []*ProtocolOrContainerUnit
}

type PCOReadingState int

const (
	ReadingID PCOReadingState = iota
	ReadingLength
	ReadingContent
)

func NewProtocolOrContainerUnit() (pcu *ProtocolOrContainerUnit) {
	pcu = &ProtocolOrContainerUnit{
		ProtocolOrContainerID: 0,
		LengthOfContents:      0,
		Contents:              []byte{},
	}
	return
}

func NewProtocolConfigurationOptions() (pco *ProtocolConfigurationOptions) {

	pco = &ProtocolConfigurationOptions{
		ProtocolOrContainerList: make([]*ProtocolOrContainerUnit, 0),
	}

	return
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) Marshal() ([]byte, error) {

	var metaInfo uint8
	var extension uint8 = 1
	var spare uint8 = 0
	var configurationProtocol uint8 = 0
	buffer := new(bytes.Buffer)

	metaInfo = (extension << 7) | (spare << 6) | (configurationProtocol)
	if err := binary.Write(buffer, binary.BigEndian, &metaInfo); err != nil {
		return nil, fmt.Errorf("Write metaInfo failed: %+v", err)
	}

	for _, containerUnit := range protocolConfigurationOptions.ProtocolOrContainerList {

		if err := binary.Write(buffer, binary.BigEndian, &containerUnit.ProtocolOrContainerID); err != nil {
			return nil, fmt.Errorf("Write protocolOrContainerID failed: %+v", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &containerUnit.LengthOfContents); err != nil {
			return nil, fmt.Errorf("Write length of contents failed: %+v", err)
		}
		if err := binary.Write(buffer, binary.BigEndian, &containerUnit.Contents); err != nil {
			return nil, fmt.Errorf("Write contents failed: %+v", err)
		}
	}
	return buffer.Bytes(), nil
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) UnMarshal(data []byte) error {
	var Buf uint8
	numOfBytes := len(data)
	byteReader := bytes.NewReader(data)
	if err := binary.Read(byteReader, binary.BigEndian, &Buf); err != nil {
		return err
	}

	numOfBytes = numOfBytes - 1
	readingState := ReadingID
	var curContainer *ProtocolOrContainerUnit

	for numOfBytes > 0 {

		switch readingState {
		case ReadingID:
			curContainer = NewProtocolOrContainerUnit()
			if err := binary.Read(byteReader, binary.BigEndian, &curContainer.ProtocolOrContainerID); err != nil {
				return err
			}
			readingState = ReadingLength
			numOfBytes = numOfBytes - 2
		case ReadingLength:
			if err := binary.Read(byteReader, binary.BigEndian, &curContainer.LengthOfContents); err != nil {
				return err
			}
			readingState = ReadingContent
			numOfBytes = numOfBytes - 1
			if curContainer.LengthOfContents == 0 {
				protocolConfigurationOptions.ProtocolOrContainerList = append(
					protocolConfigurationOptions.ProtocolOrContainerList, curContainer)
			}
		case ReadingContent:
			if curContainer.LengthOfContents > 0 {
				curContainer.Contents = make([]uint8, curContainer.LengthOfContents)
				if err := binary.Read(byteReader, binary.BigEndian, curContainer.Contents); err != nil {
					return err
				}
				protocolConfigurationOptions.ProtocolOrContainerList = append(
					protocolConfigurationOptions.ProtocolOrContainerList, curContainer)
			}
			numOfBytes = numOfBytes - int(curContainer.LengthOfContents)
			readingState = ReadingID
		}
	}
	return nil
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddDNSServerIPv4AddressRequest() {
	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.DNSServerIPv4AddressRequestUL
	protocolOrContainerUnit.LengthOfContents = 0

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddDNSServerIPv6AddressRequest() {
	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.DNSServerIPv6AddressRequestUL
	protocolOrContainerUnit.LengthOfContents = 0

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddIPAddressAllocationViaNASSignallingUL() {
	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.IPAddressAllocationViaNASSignallingUL
	protocolOrContainerUnit.LengthOfContents = 0

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddDNSServerIPv4Address(dnsIP net.IP) (err error) {

	if dnsIP.To4() == nil {
		err = fmt.Errorf("The DNS IP should be IPv4 in AddDNSServerIPv4Address!")
		return
	}
	dnsIP = dnsIP.To4()

	if len(dnsIP) != net.IPv4len {
		err = fmt.Errorf("The length of DNS IPv4 is wrong!")
		return
	}

	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.DNSServerIPv4AddressDL
	protocolOrContainerUnit.LengthOfContents = uint8(net.IPv4len)
	protocolOrContainerUnit.Contents = append(protocolOrContainerUnit.Contents, dnsIP.To4()...)

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
	return
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddPCSCFIPv4Address(pcscfIP net.IP) (err error) {
	if pcscfIP.To4() == nil {
		err = fmt.Errorf("The P-CSCF IP should be IPv4!")
		return
	}
	pcscfIP = pcscfIP.To4()

	if len(pcscfIP) != net.IPv4len {
		err = fmt.Errorf("The length of P-CSCF IP IPv4 is wrong!")
		return
	}

	protocolOrContainerUnit := NewProtocolOrContainerUnit()
	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.PCSCFIPv4AddressDL
	protocolOrContainerUnit.LengthOfContents = uint8(net.IPv4len)
	protocolOrContainerUnit.Contents = append(protocolOrContainerUnit.Contents, pcscfIP.To4()...)

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
	return
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddDNSServerIPv6Address(dnsIP net.IP) (err error) {

	if dnsIP.To16() == nil {
		err = fmt.Errorf("The DNS IP should be IPv6 in AddDNSServerIPv6Address!")
		return
	}

	if len(dnsIP) != net.IPv6len {
		err = fmt.Errorf("The length of DNS IPv6 is wrong!")
		return
	}

	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.DNSServerIPv6AddressDL
	protocolOrContainerUnit.LengthOfContents = uint8(net.IPv6len)
	protocolOrContainerUnit.Contents = append(protocolOrContainerUnit.Contents, dnsIP.To16()...)

	protocolConfigurationOptions.ProtocolOrContainerList = append(protocolConfigurationOptions.ProtocolOrContainerList,
		protocolOrContainerUnit)
	return
}

func (protocolConfigurationOptions *ProtocolConfigurationOptions) AddIPv4LinkMTU(mtu uint16) (err error) {
	protocolOrContainerUnit := NewProtocolOrContainerUnit()

	protocolOrContainerUnit.ProtocolOrContainerID = nasMessage.IPv4LinkMTUDL
	protocolOrContainerUnit.LengthOfContents = 2
	protocolOrContainerUnit.Contents =
		append(protocolOrContainerUnit.Contents, []byte{uint8(mtu >> 8), uint8(mtu & 0xff)}...)

	protocolConfigurationOptions.ProtocolOrContainerList =
		append(protocolConfigurationOptions.ProtocolOrContainerList, protocolOrContainerUnit)
	return
}
