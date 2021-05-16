package chainStruct

type Chain interface {
	Verify() bool
	Append([]byte)
	ToBytes() ([]byte, error)
	Update([]byte) error

	Array() [][]byte
	String() string
}
