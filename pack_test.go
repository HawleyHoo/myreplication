package mysql_replication_listener

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadPackTotal(t *testing.T) {
	mockBuff := []byte{0x03, 0x00, 0x00, 0x0a, 0x01, 0x02, 0x03, 0x03, 0x00, 0x00, 0x0b, 0x04, 0x05, 0x06}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, err := reader.readNextPack()

	if err != nil {
		t.Error("Got error", err)
	}

	var expectedLength uint32 = 3

	if pack.length != expectedLength {
		t.Error(
			"incorrect length",
			"expected", expectedLength,
			"got", pack.length,
		)
	}

	var expectedSequence byte = 10

	if pack.sequence != expectedSequence {
		t.Error(
			"incorrect sequence",
			"expected", expectedSequence,
			"got", pack.sequence,
		)
	}

	expectedBuff := []byte{0x01, 0x02, 0x03}

	if !reflect.DeepEqual(expectedBuff, pack.buff) {
		t.Error(
			"incorrect buff",
			"expected", expectedBuff,
			"got", pack.buff,
		)
	}

	pack, err = reader.readNextPack()

	if err != nil {
		t.Error("Got error", err)
	}

	if pack.length != expectedLength {
		t.Error(
			"incorrect length",
			"expected", expectedLength,
			"got", pack.length,
		)
	}

	expectedSequence = 11

	if pack.sequence != expectedSequence {
		t.Error(
			"incorrect sequence",
			"expected", expectedSequence,
			"got", pack.sequence,
		)
	}

	expectedBuff = []byte{0x04, 0x05, 0x06}

	if !reflect.DeepEqual(expectedBuff, pack.buff) {
		t.Error(
			"incorrect buff",
			"expected", expectedBuff,
			"got", pack.buff,
		)
	}
}

func TestReadPackByte(t *testing.T) {
	mockBuff := []byte{
		0x01, 0x00, 0x00,
		0x0a,
		0x10,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expected byte = 16
	var result byte

	err := pack.readByte(&result)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result != expected {
		t.Error(
			"Incorrect result",
			"expected", expected,
			"got", result,
		)
	}
}

func TestReadUint16(t *testing.T) {
	mockBuff := []byte{
		0x02, 0x00, 0x00,
		0x0a,
		0x1D, 0x86,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expected uint16 = 34333
	var result uint16

	err := pack.readUint16(&result)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result != expected {
		t.Error(
			"Incorrect result",
			"expected", expected,
			"got", result,
		)
	}
}

func TestReadThreeByteUint32(t *testing.T) {
	mockBuff := []byte{
		0x03, 0x00, 0x00,
		0x0a,
		0x76, 0x8A, 0x34,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expected uint32 = 3443318
	var result uint32

	err := pack.readThreeByteUint32(&result)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result != expected {
		t.Error(
			"Incorrect result",
			"expected", expected,
			"got", result,
		)
	}
}

func TestReadUint32(t *testing.T) {
	mockBuff := []byte{
		0x04, 0x00, 0x00,
		0x0a,
		0xD6, 0x00, 0x77, 0x14,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expected uint32 = 343343318
	var result uint32

	err := pack.readUint32(&result)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result != expected {
		t.Error(
			"Incorrect result",
			"expected", expected,
			"got", result,
		)
	}
}

func TestPackReadUint64(t *testing.T) {
	mockBuff := []byte{
		0x08, 0x00, 0x00,
		0x0a,
		0xC4, 0x74, 0x77, 0xCE, 0xCF, 0x11, 0x5E, 0x20,
	}

	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expected uint64 = 2332321241244333252
	var result uint64

	err := pack.readUint64(&result)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result != expected {
		t.Error(
			"Incorrect result",
			"expected", expected,
			"got", result,
		)
	}
}

func TestNilString(t *testing.T) {
	mockBuff := []byte{
		0x1C,
		0x00, 0x00, 0x0a,
		0x35, 0x2e, 0x35, 0x2e, 0x33, 0x38, 0x2d, 0x30, 0x75, 0x62, 0x75, 0x6e, 0x74, 0x75, 0x30, 0x2e, 0x31, 0x34,
		0x2e, 0x30, 0x34, 0x2e, 0x31, 0x2d, 0x6c, 0x6f, 0x67, 0x00,
	}

	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	expected := []byte("5.5.38-0ubuntu0.14.04.1-log")

	result, err := pack.readNilString()

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Error("Expected", string(expected), "got", string(result))
	}
}

func TestReadTotal(t *testing.T) {
	mockBuff := []byte{
		0x14, 0x00, 0x00,
		0x0a,
		0x10,
		0x1D, 0x86,
		0x76, 0x8A, 0x34,
		0xD6, 0x00, 0x77, 0x14,
		0xC4, 0x74, 0x77, 0xCE, 0xCF, 0x11, 0x5E, 0x20,
	}

	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	var expectedByte byte = 16
	var expected16 uint16 = 34333
	var expectedtb32 uint32 = 3443318
	var expected32 uint32 = 343343318
	var expected64 uint64 = 2332321241244333252

	var resultByte byte
	var result16 uint16
	var resulttb32 uint32
	var result32 uint32
	var result64 uint64

	err := pack.readByte(&resultByte)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if resultByte != expectedByte {
		t.Error(
			"Incorrect result",
			"expected", expectedByte,
			"got", resultByte,
		)
	}

	err = pack.readUint16(&result16)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result16 != expected16 {
		t.Error(
			"Incorrect result",
			"expected", expected16,
			"got", result16,
		)
	}

	err = pack.readThreeByteUint32(&resulttb32)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if resulttb32 != expectedtb32 {
		t.Error(
			"Incorrect result",
			"expected", expectedtb32,
			"got", resulttb32,
		)
	}

	err = pack.readUint32(&result32)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result32 != expected32 {
		t.Error(
			"Incorrect result",
			"expected", expected32,
			"got", result32,
		)
	}

	err = pack.readUint64(&result64)
	if err != nil {
		t.Error(
			"Got error", err,
		)
	}

	if result64 != expected64 {
		t.Error(
			"Incorrect result",
			"expected", expected64,
			"got", result64,
		)
	}
}

func TestWritePackUint16(t *testing.T) {
	pack := newPack()

	expected := []byte{
		0x00, 0x00, 0x00,
		0x00,
		0x30, 0x82,
	}

	var data uint16 = 33328

	err := pack.writeUInt16(data)

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, pack.Bytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestWritePackThreeByteUint32(t *testing.T) {
	pack := newPack()

	expected := []byte{
		0x00, 0x00, 0x00,
		0x00,
		0x76, 0x8A, 0x34,
	}

	var data uint32 = 3443318

	err := pack.writeThreeByteUInt32(data)

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, pack.Bytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestWritePackUint32(t *testing.T) {
	pack := newPack()

	expected := []byte{
		0x00, 0x00, 0x00,
		0x00,
		0xD6, 0x00, 0x77, 0x14,
	}

	var data uint32 = 343343318

	err := pack.writeUInt32(data)

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, pack.Bytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestWriteNilString(t *testing.T) {
	pack := newPack()

	expected := []byte{
		0x00, 0x00, 0x00,
		0x00,
		0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00,
	}

	err := pack.writeStringNil("hello")

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, pack.Bytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestPackWithlLength(t *testing.T) {
	pack := newPack()
	pack.setSequence(byte(10))

	expected := []byte{
		0x06, 0x00, 0x00,
		0x0A,
		0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00,
	}

	err := pack.writeStringNil("hello")

	if err != nil {
		t.Error("Got error", err)
	}

	if !reflect.DeepEqual(expected, pack.packBytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestPackFlush(t *testing.T) {
	mockBuff := bytes.NewBuffer([]byte{})
	packWriter := newPackWriter(mockBuff)

	pack := newPack()
	pack.setSequence(byte(10))
	pack.writeStringNil("hello")

	err := packWriter.flush(pack)

	if err != nil {
		t.Error("Got error", err)
	}

	expected := []byte{
		0x06, 0x00, 0x00,
		0x0A,
		0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00,
	}

	if !reflect.DeepEqual(expected, mockBuff.Bytes()) {
		t.Error("Expected", expected, "got", pack.Bytes())
	}
}

func TestOkPacket(t *testing.T) {
	mockBuff := []byte{
		//length
		0x07, 0x00, 0x00,
		//sequence id
		0x02,
		//code
		0x00,
		0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	if pack.isError() != nil {
		t.Error(
			"Got error", pack.isError(),
		)
	}
}

func TestOkPacketError(t *testing.T) {
	mockBuff := []byte{
		//length
		0x17, 0x00, 0x00,
		//sequence
		0x01,
		//err code
		0xff,
		//error id
		0x48, 0x04,
		//error text
		0x23, 0x48, 0x59, 0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x20, 0x75, 0x73,
		0x65, 0x64,
	}
	reader := newPackReader(bytes.NewBuffer(mockBuff))

	pack, _ := reader.readNextPack()

	errorText := "#HY000No tables used"

	err := pack.isError()

	if err == nil || err.Error() != errorText {
		t.Error(
			"incorrect err packet",
			"expected", errorText,
			"got", err.Error(),
		)
	}
}
