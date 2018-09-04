package utils

import (
	"fmt"
	"hash/adler32"
)

func DataHash(url string) string {
	// returns the hash of the url
	// could try to use https://github.com/matoous/go-nanoid
	const Size = 4
	// return crc32.ChecksumIEEE([]byte(url)) // CRC
	return fmt.Sprint(adler32.Checksum([]byte(url))) //ADLER
}
