package algorithms

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// set initial hash values

// pads message such that its length is a multiple of 512s
func padding(x string) string {
	bin_rep := ""
	for _, v := range x {
		bin_rep = fmt.Sprintf("%s%.8b", bin_rep, v)
	}
	bin_rep += "1"
	msg_len := len(x) * 8

	//find k such that adding k number of 0s makes the length of message a multiple of 512
	current_len := len(bin_rep) + 64
	var k = 0
	if current_len%512 != 0 { //if length is not already a multiple of 512
		i := current_len / 512                                    //divide length by 512 to get integer part
		k = int(math.Abs(float64(current_len - (512 * (i + 1))))) // k is then length - 512 * i

	}

	bin_rep += strings.Repeat("0", k)        //add k number of zeros
	bin_len := fmt.Sprintf("%.64b", msg_len) //convert message length to 64-bit binary
	bin_rep += bin_len
	return bin_rep
}

// parses the padded message into 512-bit blocks of 16 32-bit words
func parsing(bin_str string) [][]uint32 {
	var M [][]uint32
	N := len(bin_str) / 512
	for i := 0; i < N; i++ {

		msg_block := bin_str[i*512 : (i*512)+512] //split into 512-bit blocks
		var words []uint32
		for j := 0; j < 16; j++ {
			word := msg_block[j*32 : (j*32)+32]          //split into 32-bit word
			word_int, _ := strconv.ParseInt(word, 2, 64) //convert binary from string representation to integer
			words = append(words, uint32(word_int))

		}
		M = append(M, words)
	}
	return M

}

// rotate right
func rotr(x uint32, n int) uint32 {
	return (x >> n) | (x << (32 - n))
}

// shift right
func shr(x uint32, n int) uint32 {
	return x >> n
}

// additon then modulo 2^32
func add(x ...uint32) uint32 {
	//sum values
	var sum uint32
	for _, v := range x {
		sum += v
	}
	//2^32 is greater than 32bits, therefore must convert to int64 then back to int32
	res := uint64(sum) % uint64(math.Pow(2, 32))
	return uint32(res)
}

// sha256 functions
func Ch(x, y, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func Maj(x, y, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

func Sigma_0(x uint32) uint32 {
	return rotr(x, 2) ^ rotr(x, 13) ^ rotr(x, 22)
}

func Sigma_1(x uint32) uint32 {
	return rotr(x, 6) ^ rotr(x, 11) ^ rotr(x, 25)
}

func sigma_0(x uint32) uint32 {
	return rotr(x, 7) ^ rotr(x, 18) ^ shr(x, 3)
}

func sigma_1(x uint32) uint32 {
	return rotr(x, 17) ^ rotr(x, 19) ^ shr(x, 10)
}

// compression functions, updates the working variables
func compression(a, b, c, d, e, f, g, h, kt, wt uint32) (a_n, b_n, c_n, d_n, e_n, f_n, g_n, h_n uint32) {
	T1 := add(h, Sigma_1(e), Ch(e, f, g), kt, wt)
	T2 := add(Sigma_0(a), Maj(a, b, c))
	h_n = g
	g_n = f
	f_n = e
	e_n = add(d, T1)
	d_n = c
	c_n = b
	b_n = a
	a_n = add(T1, T2)

	return a_n, b_n, c_n, d_n, e_n, f_n, g_n, h_n

}

func Sha_256(message string) string {

	//define initial hash values and constants
	var H = []uint32{1779033703, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}

	// constants
	var K = []uint32{0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

	//preprocessing
	padded := padding(message)
	M := parsing(padded)

	//loop through each message block
	for _, m := range M {
		//define working variables
		var a, b, c, d, e, f, g, h = H[0], H[1], H[2], H[3], H[4], H[5], H[6], H[7]

		var W []uint32
		var wt uint32
		//message schedule
		for t := 0; t < 64; t++ {
			//first 16 words
			if t >= 0 && t <= 15 {
				wt = m[t]
				W = append(W, wt)

			} else {
				wt = add(sigma_1(W[t-2]), W[t-7], sigma_0(W[t-15]), W[t-16])
				W = append(W, wt)
			}
			//update working variables
			a, b, c, d, e, f, g, h = compression(a, b, c, d, e, f, g, h, K[t], wt)

		}
		//new hash values
		H[0] = add(a, H[0])
		H[1] = add(b, H[1])
		H[2] = add(c, H[2])
		H[3] = add(d, H[3])
		H[4] = add(e, H[4])
		H[5] = add(f, H[5])
		H[6] = add(g, H[6])
		H[7] = add(h, H[7])

	}

	//concert to hex and append them together
	hash := fmt.Sprintf("%.8x", H)
	hash = strings.Trim(strings.ReplaceAll(hash, " ", ""), "[]")
	return hash
}
