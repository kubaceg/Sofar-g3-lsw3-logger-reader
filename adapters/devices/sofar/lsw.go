package sofar

import (
	"encoding/binary"
	"fmt"
	"log/slog"

	"github.com/sigurn/crc16"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

type LSWRequest struct {
	serialNumber  uint
	startRegister int
	endRegister   int
}

func NewLSWRequest(serialNumber uint, startRegister int, endRegister int) LSWRequest {
	return LSWRequest{
		serialNumber:  serialNumber,
		startRegister: startRegister,
		endRegister:   endRegister,
	}
}

func (l LSWRequest) ToBytes() []byte {
	buf := make([]byte, 36)

	// preamble
	buf[0] = 0xa5
	binary.BigEndian.PutUint16(buf[1:], 0x1700)
	binary.BigEndian.PutUint16(buf[3:], 0x1045)
	buf[5] = 0x00
	buf[6] = 0x00

	// fmt.Printf("serial number: %0X\n", uint32SerialNumber)
	binary.LittleEndian.PutUint32(buf[7:], uint32(l.serialNumber))

	buf[11] = 0x02

	binary.BigEndian.PutUint16(buf[26:], 0x0103)

	binary.BigEndian.PutUint16(buf[28:], uint16(l.startRegister))
	binary.BigEndian.PutUint16(buf[30:], uint16(l.endRegister-l.startRegister+1))

	// compute crc
	table := crc16.MakeTable(crc16.CRC16_MODBUS)
	modbusCRC := crc16.Checksum(buf[26:32], table)

	// append crc
	binary.LittleEndian.PutUint16(buf[32:], modbusCRC)

	// compute & append frame crc
	buf[34] = l.checksum(buf)

	// end of frame
	buf[35] = 0x15

	return buf

}

func (l LSWRequest) String() string {
	return fmt.Sprintf("% 0X", l.ToBytes())
}

func (l LSWRequest) checksum(buf []byte) uint8 {
	var checksum uint8
	for _, b := range buf[1 : len(buf)-2] {
		checksum += b
	}
	return checksum
}

func readData(connPort ports.CommunicationPort, serialNumber uint, nameFilter func(string) bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, rr := range allRegisterRanges {
		reply, err := readRegisterRange(rr, connPort, serialNumber)
		if err != nil {
			return nil, err
		}
		for k, v := range reply {
			if nameFilter(k) {
				result[k] = v
			}
		}
	}
	return result, nil
}

func readRegisterRange(rr registerRange, connPort ports.CommunicationPort, serialNumber uint) (map[string]interface{}, error) {
	lswRequest := NewLSWRequest(serialNumber, rr.start, rr.end)

	commandBytes := lswRequest.ToBytes()

	err := connPort.Open()
	if err != nil {
		return nil, err
	}

	defer func(connPort ports.CommunicationPort) {
		if err := connPort.Close(); err != nil {
			slog.Error(fmt.Sprintf("error during connection close: %s", err))
		}
	}(connPort)

	// send the command
	_, err = connPort.Write(commandBytes)
	if err != nil {
		return nil, err
	}

	// read enough bytes
	buf := []byte{}
	for {
		b := make([]byte, 2048)
		n, err := connPort.Read(b)
		if n > 0 {
			buf = append(buf, b[:n]...)
		}
		if err != nil {
			return nil, err
		}
		if len(buf) >= 28 && len(buf) >= 28+int(buf[27]) {
			break
		}
	}
	modbusReply := buf[28 : 28+buf[27]]

	// shove the data into the reply
	reply := make(map[string]interface{})
	for _, f := range rr.replyFields {
		fieldOffset := (f.register - rr.start) * 2

		if fieldOffset > len(modbusReply)-2 {
			// skip invalid offset
			continue
		}

		switch f.valueType {
		case "U16":
			reply[f.name] = binary.BigEndian.Uint16(modbusReply[fieldOffset : fieldOffset+2])
		case "U32":
			reply[f.name] = binary.BigEndian.Uint32(modbusReply[fieldOffset : fieldOffset+4])
		case "I16":
			reply[f.name] = int16(binary.BigEndian.Uint16(modbusReply[fieldOffset : fieldOffset+2]))
		default:
		}
	}

	return reply, nil
}
