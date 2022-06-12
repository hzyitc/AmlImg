package AmlImg

import "hash/crc32"

// NOTE: Diffenent from the standard CRC32
func AmlCRC(crc uint32, p []byte) uint32 {
	table := crc32.MakeTable(0xedb88320)
	return ^crc32.Update(^crc, table, p)
}
