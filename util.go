package ipod

func BoolToByte(b bool) byte {
	if b {
		return 0x01
	}
	return 0x00
}

func ByteToBool(b byte) bool {
	return b == 0x01
}

func StringToBytes(s string) []byte {
	return append([]byte(s), 0x00)
}
