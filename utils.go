package bls

// appendAll appends multiple byte strings together
func appendAll(msgs ...[]byte) []byte {
	output := make([]byte, 0)

	for _, msg := range msgs {
		output = append(output, msg...)
	}

	return output
}
