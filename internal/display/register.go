package display

type TypeRegister int

const (
	LED TypeRegister = iota
	BUTTON
	INPUT_TEXT
	INPUT_NUM
	ARRAY_PICT
	SCREEN_NUM
)

type Register struct {
	Type   TypeRegister
	Addr   int
	Len    int
	Size   int
	Gap    int
	Toogle int
}
