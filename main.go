package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var k = 128 //kappa
var s = k + 1
var signlen = 2 * s
var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)
var lxmax = new(big.Int).Lsh(one, uint(k))         // real value space
var lxmax1 = new(big.Int).Lsh(one, uint(s))        // max value space
var lsmax = new(big.Int).Lsh(one, uint(signlen))   // sign flag
var llmax = new(big.Int).Lsh(one, uint(signlen+1)) // space

func GeneratePostiveSecret() *big.Int {
	t, _ := rand.Int(rand.Reader, lxmax)
	return t
}

func GenerateFixedSecret() *big.Int {
	x, _ := rand.Int(rand.Reader, llmax)
	return x
}

func Normalize(x *big.Int) *big.Int {
	x1 := new(big.Int).Mod(x, llmax)
	//fmt.Println("x:      ", x)
	if x1.Cmp(zero) == 1 && x1.Cmp(lsmax) == -1 {
		return x1
	} else if x1.Cmp(lsmax) == 1 || x1.Cmp(lsmax) == 0 {
		x2 := new(big.Int).Sub(llmax, x1)
		x2 = x2.Neg(x2)
		return x2
	}
	return zero
}

func ModAdd(x, y *big.Int) *big.Int {
	x1 := Normalize(x)
	y1 := Normalize(y)
	z := new(big.Int).Add(x1, y1)
	z = z.Mod(z, llmax)
	return z
}

func ModSub(x, y *big.Int) *big.Int {
	x1 := Normalize(x)
	y1 := Normalize(y)
	z := new(big.Int).Sub(x1, y1)
	z = z.Mod(z, llmax)
	return z
}

func ModMul(x, y *big.Int) *big.Int {
	x1 := Normalize(x)
	y1 := Normalize(y)
	z := new(big.Int).Mul(x1, y1)
	z = z.Mod(z, llmax)
	return z
}

func GenerateRandomShares() (*big.Int, *big.Int, *big.Int) {
	x := GenerateFixedSecret()
	x1 := GenerateFixedSecret()
	x2 := ModSub(x, x1)
	x2 = x2.Mod(x2, llmax)
	return x, x1, x2
}

func GeneratePostiveRandomShares() (*big.Int, *big.Int, *big.Int) {
	x := GeneratePostiveSecret()
	x1 := GenerateFixedSecret()
	x2 := ModSub(x, x1)
	return x, x1, x2
}

func GenerateRandomValueShares() (*big.Int, *big.Int, *big.Int) {
	x := GeneratePostiveSecret()
	r, _ := rand.Int(rand.Reader, two)
	//fmt.Println(r)
	if r.Cmp(one) == 0 {
		x = x.Neg(x)
	}
	x1 := GenerateFixedSecret()
	x2 := ModSub(x, x1)
	x2 = x2.Mod(x2, llmax)
	return x, x1, x2
}

// [x],[y] ===> [xy]
func SecMul(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	// offline
	a, a1, a2 := GenerateRandomShares()
	b, b1, b2 := GenerateRandomShares()
	c := ModMul(a, b)
	c1 := GenerateFixedSecret()
	c2 := ModSub(c, c1)

	//online
	e1 := ModSub(x1, a1)
	f1 := ModSub(y1, b1)
	e2 := ModSub(x2, a2)
	f2 := ModSub(y2, b2)

	e := ModAdd(e1, e2)
	f := ModAdd(f1, f2)

	res11 := ModMul(b1, e)
	res12 := ModMul(a1, f)
	res13 := ModMul(e, f)
	res1 := ModAdd(c1, res11)
	res1 = ModAdd(res1, res12)
	res1 = ModAdd(res1, res13)

	res21 := ModMul(b2, e)
	res22 := ModMul(a2, f)
	res2 := ModAdd(c2, res21)
	res2 = ModAdd(res2, res22)
	return res1, res2
}

func SecCmp(x1, y1, x2, y2 *big.Int) int {
	_, t1, t2 := GeneratePostiveRandomShares()
	//fmt.Println("t:", t)
	diff1 := ModSub(x1, y1)
	diff2 := ModSub(x2, y2)
	cx1, cx2 := SecMul(diff1, t1, diff2, t2)
	cx := ModAdd(cx1, cx2)
	cx = Normalize(cx)
	if cx.Cmp(zero) == 0 {
		return 0
	} else if cx.Cmp(zero) == 1 {
		return 1
	}
	return -1
}

func main() {
	count := 0
	for i := 0; i < 10000; i++ {
		x, x1, x2 := GenerateRandomValueShares()
		y, y1, y2 := GenerateRandomValueShares()
		x = Normalize(x)
		//fmt.Println(x)
		y = Normalize(y)
		//fmt.Println(y)
		sign1 := -1
		if x.Cmp(y) == 0 {
			sign1 = 0
		} else if x.Cmp(y) == 1 {
			sign1 = 1
		}
		//fmt.Println(sign1)
		sign2 := SecCmp(x1, y1, x2, y2)
		//fmt.Println(sign2)
		if sign1 != sign2 {
			count = count + 1
		}
	}
	fmt.Println(count)
}
