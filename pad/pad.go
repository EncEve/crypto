// Use of this source code is governed by a license
// that can be found in the LICENSE file.

// Package pad implements some padding schemes
// for block ciphers.
package pad

import (
	cryptorand "crypto/rand"
	"errors"
	"io"
)

var badPadErr = errors.New("bad padding bytes")
var notMulOfBlockErr = errors.New("src is not a multiply of the padding blocksize")

// The Padding interface represents a padding scheme.
type Padding interface {

	// BlockSize returns the block size of the padding.
	BlockSize() int

	// Returns the overhead, the padding will cause
	// by padding the given byte slice. The overhead
	// will always be between 1 and BlockSize() inclusively.
	Overhead(src []byte) int

	// Pads the last (may incomplete) block of the src slice
	// to a padded and complete block, appends the padding bytes
	// to the src slice and returns this slice.
	// The length of the returned slice is len(src) + Overhead(src)
	Pad(src []byte) []byte

	// Takes a slice and tries to remove the padding bytes
	// form the last block. Therefore the length of the
	// src argument must be a multiply of the blocksize.
	// If the returned error is nil, the padding could be
	// removed successfully. The returned slice holds the
	// unpadded src bytes.
	Unpad(src []byte) ([]byte, error)
}

// NewX923 returns a new pad.Padding implementing the ANSI X.923 scheme.
// Only block sizes between 1 and 255 are valid.
func NewX923(blocksize int) Padding {
	if blocksize < 1 || blocksize > 255 {
		panic("illegal blocksize - size must between 0 and 256")
	}
	pad := x923Padding(blocksize)
	return pad
}

// NewPKCS7 returns a new pad.Padding implementing the PKCS 7 scheme.
// Only block sizes between 1 and 255 are valid.
func NewPKCS7(blocksize int) Padding {
	if blocksize < 1 || blocksize > 255 {
		panic("illegal blocksize - size must between 0 and 256")
	}
	pad := pkcs7Padding(blocksize)
	return pad
}

// NewISO10126 returns a new pad.Padding, which uses the padding scheme
// described in ISO 10126. The padding bytes are taken
// form the given rand argument. If rand is nil, crypto/rand will be used.
// Only block sizes between 1 and 255 are valid.
func NewISO10126(blocksize int, rand io.Reader) Padding {
	if blocksize < 1 || blocksize > 255 {
		panic("illegal blocksize - size must between 0 and 256")
	}
	pad := new(isoPadding)
	pad.blocksize = blocksize
	if rand == nil {
		pad.random = cryptorand.Reader
	} else {
		pad.random = rand
	}
	return pad
}

// Returns the overhead for a given slice with a
// specific block size.
func overhead(blocksize int, src []byte) int {
	return blocksize - (len(src) % blocksize)
}
