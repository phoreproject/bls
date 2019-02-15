package bls

// BitIterator is an iterator through bits.
type BitIterator struct {
	arr []uint64
	n   uint
}

// NewBitIterator creates a new bit iterator given an array of ints.
func NewBitIterator(arr []uint64) BitIterator {
	return BitIterator{arr, uint(len(arr) * 64)}
}

// Next returns the next bit in the bit iterator with the
// second return value as true when finished.
func (bi *BitIterator) Next() (bool, bool) {
	if bi.n == 0 {
		return false, true
	}
	bi.n--
	part := bi.n / 64
	bit := bi.n - (part * 64)
	return bi.arr[part]&(1<<bit) > 0, false
}
