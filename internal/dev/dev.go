package dev

import (
	"encoding/binary"
	"io"
	"os"
)

// CopyFrom copies all data from device to writer.
func CopyFrom(file *os.File, device string) (size int64, err error) {

	// Take a source device.
	dev, err := os.OpenFile(device, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer dev.Close()

	// Reading a real size of data and copy data to writer.
	if err = binary.Read(dev, binary.LittleEndian, &size); err != nil {
		return 0, err
	}
	return io.CopyN(file, dev, size)
}

// CopyTo copies all data from reader to device.
func CopyTo(device string, file *os.File) (size int64, err error) {

	// Take a target device.
	dev, err := os.OpenFile(device, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer dev.Close()

	// Take a file size and write it to device.
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	err = binary.Write(dev, binary.LittleEndian, stat.Size())
	if err != nil {
		return 0, err
	}

	// Write data to device.
	return io.CopyN(dev, file, stat.Size())
}
