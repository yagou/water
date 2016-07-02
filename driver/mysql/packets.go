package mysql

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (m *mysql) readPacket() ([]byte, error) {
	// 获取包头
	data, err := m.buf.readNext(4)

	if err != nil {
		errLog.Print(err.Error())
		m.Close()
		return nil, errBadConn
	}

	// 获取包长度[24 bit]
	pktLen := int(uint32(data[0]) | uint32(data[1]<<8) | uint32(data[2]<<16))

	if pktLen < 1 {
		errLog.Print(errMalformPkt.Error())
		m.Close()
		return nil, errBadConn
	}

	// Check Packet Sync [8 bit]
	if data[3] != m.sequence {
		if data[3] > m.sequence {
			return nil, errPktSyncMul
		} else {
			return nil, errPktSync
		}
	}
	m.sequence++
	// 获取包数据长度
	if data, err = m.buf.readNext(pktLen); err == nil {
		if pktLen < maxPacketSize {
			return data, nil
		}

		buf := make([]byte, len(data))
		copy(buf, data)

		data, err = m.readPacket()
		if err == nil {
			return append(buf, data...), nil
		}
	}

	m.Close()
	errLog.Print(err.Error())
	return nil, errBadConn
}

func (m *mysql) readInitPacket() ([]byte, error) {
	data, err := m.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == iERR {
		return nil, m.handleErrorPacket(data)
	}

	if data[0] < minProtocolVersion {
		return nil, fmt.Errorf(
			"Unsupported MySQL Protocol Version %d. Protocol Version %d or higher is required",
			data[0],
			minProtocolVersion,
		)
	}
	// server version [null terminated string]
	// connection id [4 bytes]
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4
	// first part of the password cipher [8 bytes]
	cipher := data[pos : pos+8]
	// (filler) always 0x00 [1 byte]
	pos += 8 + 1

	// capability flags (lower 2 bytes) [2 bytes]
	m.flags = clientFlag(binary.LittleEndian.Uint16(data[pos : pos+2]))

	if m.flags&clientProtocol41 == 0 {
		return nil, errOldProtocol
	}

	pos += 2
	if len(data) > pos {
		// character set [1 byte]
		// status flags [2 bytes]
		// capability flags (upper 2 bytes) [2 bytes]
		// length of auth-plugin-data [1 byte]
		// reserved (all [00]) [10 bytes]
		pos += 1 + 2 + 2 + 1 + 10

		// second part of the password cipher [12? bytes]
		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// which is not documented but seems to work.
		cipher = append(cipher, data[pos:pos+12]...)

		// TODO: Verify string termination
		// EOF if version (>= 5.5.7 and < 5.5.10) or (>= 5.6.0 and < 5.6.2)
		// \NUL otherwise
		//
		//if data[len(data)-1] == 0 {
		//	return
		//}
		//return errMalformPkt
	}

	return cipher, nil

}

// Error Packet
// http://dev.mysql.com/doc/internals/en/generic-response-packets.html#packet-ERR_Packet
func (m *mysql) handleErrorPacket(data []byte) error {
	if data[0] != iERR {
		return errMalformPkt
	}

	// 0xff [1 byte]

	// Error Number [16 bit uint]
	errno := binary.LittleEndian.Uint16(data[1:3])

	pos := 3

	// SQL State [optional: # + 5bytes string]
	if data[3] == 0x23 {
		//sqlstate := string(data[4 : 4+5])
		pos = 9
	}

	// Error Message [string]
	return &MySQLError{
		Number:  errno,
		Message: string(data[pos:]),
	}
}
