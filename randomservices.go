package ut

import "math/rand"

type RandomServices struct{}

// RandomArray returns a deterministically seeded random array of the given length
func (rs *RandomServices) RandomArray(i, length int) []byte {
	source := rand.NewSource(int64(i))
	r := rand.New(source)
	b := make([]byte, length)
	for n := 0; n < length; n++ {
		b[n] = byte(r.Intn(256))
	}
	return b
}
