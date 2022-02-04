package chunk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

func Write(addressToAppearances map[string][]AppearanceRecord) (*bytes.Buffer, error) {
	// First, prepare some placeholders for our data
	addressTable := make([]AddressRecord, 0, len(addressToAppearances))
	appearanceCount := 0
	// We need to write appearance table after address table, but we can build it first,
	// so we will use another buffer to store it as bytes
	appearanceTableBuf := bytes.Buffer{}

	for address, appearances := range addressToAppearances {
		// Convert string with an address into bytes
		addressHex, err := hex.DecodeString(address[2:])
		if err != nil {
			return nil, err
		}
		ethAddr := EthAddress{}
		copy(ethAddr[:], addressHex)

		// Append record to address table. Note that StartRecord is just total
		// number of previous records (counting from 0)
		addressTable = append(addressTable, AddressRecord{
			Address:         ethAddr,
			StartRecord:     uint32(appearanceCount),
			NumberOfRecords: uint32(len(appearances)),
		})
		appearanceCount += len(appearances)

		// Write appearance data into the "placeholder" buffer
		binary.Write(&appearanceTableBuf, binary.LittleEndian, &appearances)
	}

	// This is our main buffer. We will write all chunk contents there
	buf := bytes.Buffer{}
	header := Header{
		Magic:               MagicNumber,
		Hash:                EthHash{},
		NumberOfAddresses:   uint32(len(addressToAppearances)),
		NumberOfAppearances: uint32(appearanceCount),
	}

	// Write header
	err := binary.Write(&buf, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	// Write address table
	err = binary.Write(&buf, binary.LittleEndian, &addressTable)
	if err != nil {
		return nil, err
	}
	if buf.Len() != (44 + 84) {
		panic("oh no")
	}

	// Write appearance table
	_, err = appearanceTableBuf.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
