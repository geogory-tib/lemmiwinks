package util

import (
	"bytes"
	"io"
	"log"
	"runtime"
)

// like io.ReadAll but completes on the end of the Json Object
// this does not account for nested objects
func ReadFullJson(reader io.Reader, buffer []byte) (int, error) {
	temp_buffer := make([]byte, 512)
	total_bytes := 0
	for {
		n, err := reader.Read(temp_buffer)
		if err != nil {
			return total_bytes, err
		}
		if total_bytes+n > len(buffer) {
			return total_bytes, io.ErrShortBuffer
		}

		copy(buffer[total_bytes:(total_bytes+n)], temp_buffer[:n])
		total_bytes += n
		if bytes.Contains(temp_buffer, []byte("}")) {
			return total_bytes, err
		}
	}
}

func Todo() {
	pc, file, line, _ := runtime.Caller(1)
	log.Panicf("This has yet to be implemented yet: PC: %x, FILE: %s, LINE: %d\n", pc, file, line)
}
