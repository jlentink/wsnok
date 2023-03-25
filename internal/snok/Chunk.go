package snok

type Chunk struct {
	Index  int
	Url    string
	Offset int64
	Size   int64
	Data   []byte
}

func (c *Chunk) Write(b []byte) (n int, err error) {
	c.Data = append(c.Data, b...)
	n = len(b)
	err = nil
	return
}
