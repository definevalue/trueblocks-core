package chunk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestWrite(t *testing.T) {
	a2a := map[string][]AppearanceRecord{
		"0x0780313e308944f2985e74d480de983222b87c86": {
			{
				BlockNumber:   uint32(14134380),
				TransactionId: uint32(21),
			},
			{
				BlockNumber:   uint32(14134380),
				TransactionId: uint32(35),
			},
		},
		"0xb759c2fb021bb84f016ac63c8faba64f72b517d3": {
			{
				BlockNumber:   uint32(14134380),
				TransactionId: uint32(155),
			},
		},
		"0xbc0a8358507fd406fa97ec82aeb6fa057e9603e9": {
			{
				BlockNumber:   uint32(14134380),
				TransactionId: uint32(226),
			},
		},
	}

	addressCount := 0
	appearanceCount := 0
	for _, appearances := range a2a {
		addressCount++
		appearanceCount += len(appearances)
	}

	result, err := Write(a2a)
	if err != nil {
		t.Fatal("write error:", err)
	}

	// f, err := os.Create("/tmp/chunk.bin")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// result.WriteTo(f)
	// f.Close()

	// readChunk := func(dest interface{}) error {
	// 	return binary.Read(result, binary.LittleEndian, dest)
	// }
	// header := Header{}
	// err = readChunk(&header)
	// if err != nil {
	// 	t.Fatal("header read error:", err)
	// }

	// if header.NumberOfAddresses != 3 {
	// 	t.Fatalf("wrong number of addresses %d", header.NumberOfAddresses)
	// }
	// if header.NumberOfAppearances != 4 {
	// 	t.Fatalf("wrong number of appearances %d", header.NumberOfAppearances)
	// }

	// firstAddress := AddressRecord{}

	// err = readChunk(&firstAddress)
	// if err != nil {
	// 	t.Fatal("address read error:", err)
	// }

	// addressBytes := []byte("0x0780313e308944f2985e74d480de983222b87c86")
	// if !bytes.Equal(firstAddress.Address[:], addressBytes[2:]) {
	// 	t.Fatal("wrong first address", firstAddress.Address)
	// }

	// if firstAddress.NumberOfRecords != 2 {
	// 	t.Fatal("wrong first number of record", firstAddress.NumberOfRecords)
	// }

	rawHeader := make([]byte, HeaderSize)
	rawAddressTable := make([]byte, addressCount*AddressSize)
	rawAppearanceTable := make([]byte, appearanceCount*AppearanceSize)

	_, err = result.Read(rawHeader)
	if err != nil {
		t.Fatal(err)
	}
	_, err = result.Read(rawAddressTable)
	if err != nil {
		t.Fatal(err)
	}
	_, err = result.Read(rawAppearanceTable)
	if err != nil {
		t.Fatal(err)
	}

	readBytes := func(slice []byte, dest interface{}) error {
		r := bytes.NewReader(slice)
		return binary.Read(r, binary.LittleEndian, dest)
	}

	header := Header{}
	err = readBytes(rawHeader, &header)
	if err != nil {
		t.Fatal("reading header:", err)
	}
	if header.NumberOfAddresses != 3 {
		t.Fatalf("wrong number of addresses %d", header.NumberOfAddresses)
	}
	if header.NumberOfAppearances != 4 {
		t.Fatalf("wrong number of appearances %d", header.NumberOfAppearances)
	}

	addresses := [3]AddressRecord{}
	err = readBytes(rawAddressTable, &addresses)
	if err != nil {
		t.Fatal("reading address table:", err)
	}
	if len(addresses) == 0 {
		t.Fatal("empty address table", len(rawAddressTable))
	}
	addressBytes, err := hex.DecodeString("0x0780313e308944f2985e74d480de983222b87c86"[2:])
	t.Fatal(addressBytes)
	if err != nil {
		t.Fatal("converting address:", err)
	}
	readAddressAsBytes := []byte{}
	copy(readAddressAsBytes, addresses[0].Address[:])
	if !bytes.Equal(addresses[0].Address[:], addressBytes) {
		t.Fatal("wrong first address", addresses[0].Address)
	}

	if addresses[0].NumberOfRecords != 2 {
		t.Fatal("wrong first number of record", addresses[0].NumberOfRecords)
	}
}
