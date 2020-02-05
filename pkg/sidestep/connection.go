package sidestep

// Connection implements io.Reader and io.Writer
type Connection interface {
	Write(data []byte) (n int, err error)
	Read(data []byte) (n int, err error)
	Close() error
}
