package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"time"
)

var A, B, C, D int        //inputted coefficients before type casting to complex128
var a, b, c, d complex128 // coefficients  for f(z)
var err error
var ps = []*point{}       // the array of all input/output points
var r1, r2, r3 complex128 // the roots

const p1024 = 0.01953125 // the increment between points; 20/1024
var purple = color.RGBA{R: 75, B: 130, A: 255}
var darkBlue = color.RGBA{R: 36, G: 74, B: 150, A: 255}
var greenBlue = color.RGBA{R: 95, G: 158, B: 160, A: 255}

type point struct {
	out     complex128 //final output value
	nearest int        //nearest root, represented as a 0,1,2 or -1 if out = NaN
}

func main() {
	askCoefS()
	// fun test case: 1,-1,0,2 = -1+0i, 1+i, 1-i;
	start := time.Now()
	var results = map[complex128]int{}
	for j := -512.0; j < 512; j++ {
		for k := -512.0; k < 512; k++ {
			output := approx(complex(j*p1024, k*p1024))
			ps = append(ps, &point{output, -1})
			results[output] += 1
		}
	}

	statCrunch(results)

	m := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	var colour color.RGBA
	for i := range ps {
		if ps[i].nearest == 0 {
			colour = purple
		} else if ps[i].nearest == 1 {
			colour = darkBlue
		} else {
			colour = greenBlue
		}
		m.Set(int(i/1024), i%1024, colour)
	}
	out, err := os.Create("./images/" + ctoa(a) + "_" + ctoa(b) + "_" + ctoa(c) + "_" + ctoa(d) + ".png")
	e(err)
	err = png.Encode(out, m)
	e(err)
	fmt.Println(time.Now().Sub(start))
}

func f(z complex128) complex128 {
	return a*cmplx.Pow(z, 3) + b*cmplx.Pow(z, 2) + c*z + d
}
func fP(z complex128) complex128 {
	return 3*a*cmplx.Pow(z, 2) + 2*b*z + c
}

func approx(z complex128) complex128 { // iterates over the input points; returns final output value
	for range make([]int, 20) {
		z = z - (f(z) / fP(z))
	}
	return complex(math.Round(100*real(z))/100, math.Round(100*imag(z))/100) // rounding to 2 decimal places
}
func nearRoot(p *point) { //finds which root is nearest to a points output value for colouring purposes
	if p.out == cmplx.NaN() {
		p.nearest = -1
		return
	}
	distFromR1 := math.Hypot(real(r1)-real(p.out), imag(r1)-imag(p.out))
	distFromR2 := math.Hypot(real(r2)-real(p.out), imag(r2)-imag(p.out))
	distFromR3 := math.Hypot(real(r3)-real(p.out), imag(r3)-imag(p.out))

	nearest := math.Min(math.Min(distFromR1, distFromR2), distFromR3)
	if nearest == distFromR1 {
		p.nearest = 0
	} else if nearest == distFromR2 {
		p.nearest = 1
	} else {
		p.nearest = 2
	}
}
func askCoefS() { //record coefficients
	fmt.Println("ax^3 + bx^2 + cx + d")
	fmt.Println("Enter a non-zero value for A and any value for B, C and D")
	_, err = fmt.Scanln(&A, &B, &C, &D)
	e(err)
	a, b, c, d = itoc(A), itoc(B), itoc(C), itoc(D)
}
func statCrunch(results map[complex128]int) { // finds the roots from data and runs color assigning function
	var res = results
	var topThree = make([]complex128, 3)
	var topValue = make([]int, 3)
	var val int
	for i := 0; i < 3; i++ {
		for j, jval := range res {
			if jval > val {
				topThree[i] = j
				val = jval
			}
		}
		topValue[i] = val
		res[topThree[i]] = 0
		val = 0
	}
	r1, r2, r3 = topThree[0], topThree[1], topThree[2]
	fmt.Println(topThree)
	for i := range ps {
		nearRoot(ps[i])
	}
}

func itoc(x int) complex128 { // int to complex
	return complex(float64(x), 0)
}
func ctoa(z complex128) string { // complex to string (only the real part)
	return strconv.Itoa(int(real(z)))
}
func e(err error) { //error handling lol
	if err != nil {
		panic(err)
	}
}
