package models

type Format int

func (Format) PDF() Format {
	return 1
}

func (Format) TXT() Format {
	return 2
}
