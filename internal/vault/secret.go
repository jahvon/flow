package vault

type Secret string

func (s Secret) ObfuscatedString() string {
	if s.Empty() {
		return ""
	}
	return "********"
}

func (s Secret) String() string {
	return s.ObfuscatedString()
}

func (s Secret) PlainTextString() string {
	return string(s)
}

func (s Secret) Empty() bool {
	return string(s) == ""
}
