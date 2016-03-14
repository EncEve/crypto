// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package blake2s

import (
	"encoding/hex"
	"testing"
)

type testVector struct {
	p        *Params
	src, exp string
}

func decodeHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		if t == nil {
			panic(err)
		}
		t.Fatal(err)
	}
	return b
}

var vectors []testVector = []testVector{
	// Test vectors from https://blake2.net/blake2s-test.txt
	testVector{
		// without explicit hash size (check if default is used)
		p:   &Params{Key: decodeHex(nil, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")},
		src: hex.EncodeToString([]byte("")),
		exp: "48a8997da407876b3d79c0d92325ad3b89cbb754d86ab71aee047ad345fd2c49",
	},
	testVector{
		p: &Params{HashSize: 32,
			Key: decodeHex(nil, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")},
		src: "00",
		exp: "40d15fee7c328830166ac3f918650f807e7e01e177258cdc0a39b11f598066f1",
	},
	testVector{
		p: &Params{HashSize: 32,
			Key: decodeHex(nil, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")},
		src: "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f" +
			"202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40414" +
			"2434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f6061626364" +
			"65666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868" +
			"788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9" +
			"aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbc" +
			"ccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedee" +
			"eff0f1f2f3f4f5f6f7f8f9fafbfcfdfe",
		exp: "3fb735061abc519dfe979e54c1ee5bfad0a9d858b3315bad34bde999efd724dd",
	},
}

func testSingleVector(t *testing.T, i int) {
	vec := &vectors[i]
	h, err := New(vec.p)
	if err != nil {
		t.Fatal(err)
	}
	srcBin, err := hex.DecodeString(vec.src)
	if err != nil {
		t.Fatal(err)
	}

	h.Write(srcBin)
	sum := h.Sum(nil)
	expSum, _ := hex.DecodeString(vec.exp)
	if len(sum) != len(expSum) {
		t.Fatalf("Test vector %d : Hash size does not match expected - found %d expected %d", i, len(sum), len(expSum))
	}
	for j := range sum {
		if sum[j] != expSum[j] {
			t.Fatalf("Test vector %d : Hash does not match:\nFound:    %v\nExpected: %v", i, hex.EncodeToString(sum), hex.EncodeToString(expSum))
		}
	}
}

func TestBlake2s(t *testing.T) {
	for i := range vectors {
		testSingleVector(t, i)
	}
}