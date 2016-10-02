package upx

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

func IsElfUpxed(appName string) (bool, error) {
	f, err := os.Open(appName)
	if err != nil {
		return false, err
	}

	defer f.Close()
	_elf, err := elf.NewFile(f)

	// Read and decode ELF identifier
	var ident [16]uint8
	f.ReadAt(ident[0:], 0)

	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return false, fmt.Errorf("Bad magic number at %d\n", ident[0:4])
	}

	switch _elf.Class.String() {
	case "ELFCLASS64":
		var hdr elf.Header64
		f.Seek(0, os.SEEK_SET)
		if err := binary.Read(f, _elf.ByteOrder, hdr); err != nil {
			return false, err
		}
		_elf.Progs[0].Flags.String()
		f.Seek(int64(hdr.Phoff)+int64(hdr.Phentsize)*int64(hdr.Phnum), os.SEEK_SET)

	case "ELFCLASS32":
		var hdr elf.Header32
		f.Seek(0, os.SEEK_SET)
		if err := binary.Read(f, _elf.ByteOrder, hdr); err != nil {
			return false, err
		}
		f.Seek(int64(hdr.Phoff)+int64(hdr.Phentsize)*int64(hdr.Phnum), os.SEEK_SET)

	default:
		return false, fmt.Errorf("unsupport class", _elf.Class.String())
	}
	var upxMagic [8]byte
	if _, err := f.Read(upxMagic[0:]); err != nil {
		return false, err
	}

	return string(upxMagic[4:]) == "UPX!", nil
}
