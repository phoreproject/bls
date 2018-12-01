package bls

// MACWithCarry performs the operation a + b * c + carry and returns
// the result and sets carry to the high bits.
func MACWithCarry(a uint64, b uint64, c uint64, carry *uint64) uint64

// SubWithBorrow performs the operation a - b + borrow and returns
// the result and sets carry to the high bits.
func SubWithBorrow(a uint64, b uint64, borrow *uint64) uint64

// AddWithCarry performs the operation a + b + carry and returns
// the result and sets carry to the high bits.
func AddWithCarry(a uint64, b uint64, carry *uint64) uint64
