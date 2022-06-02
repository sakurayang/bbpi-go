package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	floatOne     = big.NewFloat(1)
	floatTwo     = big.NewFloat(2)
	floatFour    = big.NewFloat(4)
	// floatFive    = big.NewFloat(5)
	// floatSix     = big.NewFloat(6)
	// floatEight   = big.NewFloat(8)
	floatSixteen = big.NewFloat(16)
	// intOne = big.NewInt(1)
	// intFour = big.NewInt(4)
	// intFive = big.NewInt(5)
	// intSix   = big.NewInt(6)
	intEight      = big.NewInt(8)
	intSixteen       = big.NewInt(16)
	prec        uint = 0
	programTime time.Time
	out2file    = false
)

type result struct {
	Range struct {
		Start uint32 `json:"start"`
		End   uint32 `json:"end"`
	} `json:"range"`
	EnableMulti bool   `json:"enable_multi"`
	Pi          string `json:"pi"`
	TimeUse     string `json:"time_use"`
}

type piChunk struct {
	order uint32
	value string
}

func log(res result) {
	if out2file {
		filename := fmt.Sprintf("pi_%d.json", programTime.UnixMicro())
		fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		resJson, _ := json.Marshal(res)
		_, err = fd.Write(resJson)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Printf("范围: %d - %d \n", res.Range.Start, res.Range.End)
		fmt.Printf("启用多进程：%v \n", res.EnableMulti)
		fmt.Println("pi: ", res.Pi)
		fmt.Printf("用时: %v\n", res.TimeUse)
	}
}

// fpart 返回 x 的小数部分。
// 返回的 y 符合 (|x| + y) % 1 == 0 ，并且 y > 0。
func fpart(x *big.Float) *big.Float {
	a := new(big.Float).Copy(x)
	// 向下取整
	aInt, _ := a.Int(nil)
	aSign := a.Sign()
	if aSign < 0 {
		a.Add(a, new(big.Float).SetInt(aInt.Neg(aInt)))
		return a.Add(floatOne, a)
	} else if aSign > 0 {
		return a.Sub(a, new(big.Float).SetInt(aInt))
	} else {
		return a
	}
}

func bbp(n uint32, j int64, mul *big.Float) *big.Float {
	s := big.NewFloat(0).SetPrec(prec)
	// sum(16^(n-k) mod (8k+1) / (8k+1)), from 0 to n
	k8 := big.NewInt(j)
	for k := uint32(0); k <= n; k++ {
		nk := big.NewInt(int64(n - k))
		// 16^(n-k) mod (8k+1)
		a := new(big.Int).Exp(intSixteen, nk, k8)
		// / (8k+1)
		b := new(big.Float).SetPrec(prec).SetInt(a)
		b.Quo(b, new(big.Float).SetInt(k8))

		s.Add(s, b)
		k8.Add(k8, intEight)
	}
	// fmt.Println(s)

	//sum(16^(n-k) / (8k+1)), from n+1 to inf
	num := big.NewFloat(1/16)
	for k := int64(n + 1); k < int64((n+1)*2+uint32(prec)); k++ {
		frac := new(big.Float).SetPrec(prec).Quo(num, new(big.Float).SetInt(k8))
		s.Add(s, frac)
		num.Quo(num, floatSixteen)
		k8.Add(k8, intEight)
	}
	//fmt.Println(s)
	s.Mul(mul, fpart(s))
	return fpart(s)
}

/*
// bit take a formula named bbp
// BBP - https://www.wikiwand.com/en/Bailey%E2%80%93Borwein%E2%80%93Plouffe_formula
// $$pi = \sum_{k=0}^{\infty} [ \frac{1}{16^k}(\frac{4}{8k+1} - \frac{2}{8k+4}) - \frac{1}{8k+5} - \frac{1}{8k+6}]$$
// pi = sum ( 1/16^k * (4/(8k+1) - 2/(8k+4) - 1/(8k+5) - 1/(8k+6) )), k = 0 to inf
func bit(n uint32) *big.Float {
	k := big.NewFloat(float64(n))

	var k8 *big.Float
	if n == 0 {
		k8 = big.NewFloat(float64(0))
	}
	k8 = new(big.Float).Mul(k, floatEight)
	kp16 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(16), big.NewInt(int64(n)), nil))
	// kp16 = kp16.SetPrec(prec).Quo(one, kp16)

	k8_1 := new(big.Float).Add(k8, floatOne)
	k8_4 := new(big.Float).Add(k8, floatFour)
	k8_5 := new(big.Float).Add(k8, floatFive)
	k8_6 := new(big.Float).Add(k8, floatSix)
	//  4.0 / (k8 + 1.0)
	part1 := new(big.Float).Quo(floatFour, k8_1).SetPrec(prec)
	// 2.0 / (k8 + 4.0)
	part2 := new(big.Float).Quo(floatTwo, k8_4).SetPrec(prec)
	// 1.0 / (k8 + 5.0)
	part3 := new(big.Float).Quo(floatOne, k8_5).SetPrec(prec)
	// 1.0 / (k8 + 6.0)
	part4 := new(big.Float).Quo(floatOne, k8_6).SetPrec(prec)

	a := new(big.Float).Sub(part1, part2).SetPrec(prec)
	b := new(big.Float).Add(part3, part4).SetPrec(prec)
	c := a.Sub(a, b)
	c = c.Quo(c, kp16)
	//c.Sub(c, new(big.Float).SetInt(i))
	fmt.Println(c.String())
	return c
}
*/

func bit(n uint32) string{
	p1 := bbp(n, 1, floatFour)
	p2 := bbp(n, 4, floatTwo)
	p3 := bbp(n, 5, floatOne)
	p4 := bbp(n, 6, floatOne)

	// a - b - c -d === (a-b) - (c+d)
	pi := new(big.Float).SetPrec(prec).Sub(p1.Sub(p1, p2), p3.Add(p3, p4))
	pi = pi.Mul(floatSixteen, fpart(pi))
	pInt, _ := pi.Int(nil)
	return  fmt.Sprintf("%x", pInt)
}

func chunk(start, end uint32) string {
	pi := ""
	for k := start; k < end; k++ {
		pi += bit(k)
	}
	return pi
}

func multiProcess(end uint32, start uint32) string {
	pi := ""
	cores := uint32(runtime.NumCPU())
	runtime.GOMAXPROCS(int(cores))
	c := make(chan *piChunk, cores)
	finish := make(chan bool, cores)

	bitsPreChunkFloat := float64(end) / float64(cores)
	bitsPreChunk := uint32(math.Floor(bitsPreChunkFloat))
	last := end - bitsPreChunk*cores

	if end-start < cores {
		bitsPreChunk = end - start
		last = bitsPreChunk
		cores = 1
	}

	for k := uint32(0); k < cores; k++ {
		go func(k uint32) {
			t := ""
			if k != cores {
				t = chunk(k*bitsPreChunk+start, (k+1)*bitsPreChunk-1+start)
			} else {
				t = chunk(k*bitsPreChunk+start, k*bitsPreChunk+last+start)
			}
			c <- &piChunk{
				order: k,
				value: t,
			}
			finish <- true
		}(k)
	}
	piSortMap := map[uint32]string{}
	for i := uint32(0); i < cores; i++ {
		<-finish
		t := <-c

		piSortMap[t.order] = t.value
	}

	close(c)
	close(finish)

	for i := uint32(0); i < cores; i++ {
		pi += piSortMap[i]
	}

	return pi
}

func main() {
	// ========= init =============
	if len(os.Args) < 4 {
		fmt.Println("命令格式： pi [起始位] [结束位] [是否启用多线程模式] <是否写入文件>")
		fmt.Println("起始位：数字，最小为0")
		fmt.Println("结束位：数字，最大为4294967295")
		fmt.Println("是否启用多线程：布尔值，true 或 false")
		fmt.Println("是否写入到文件：布尔值，默认为 false，启用后会在同目录下生成一个以时间戳命名的 json 文件")
		os.Exit(1)
	}

	prec = uint(640)
	programTime = time.Now()

	arg1 := os.Args[1]
	arg2 := os.Args[2]
	arg3 := os.Args[3]

	argInt, _ := strconv.Atoi(arg1)
	start := uint32(argInt)

	argInt, _ = strconv.Atoi(arg2)
	end := uint32(argInt)

	enableMulti, _ := strconv.ParseBool(arg3)

	if len(os.Args) >= 5 {
		out2file, _ = strconv.ParseBool(os.Args[4])
	}

	// ========= calc =============
	begin := time.Now()

	pi := ""

	if enableMulti {
		pi = multiProcess(end, start)
	} else {
		pi = chunk(start, end)
	}

	over := time.Now()

	// ========= end =============
	res := result{
		Range: struct {
			Start uint32 `json:"start"`
			End   uint32 `json:"end"`
		}{
			Start: start,
			End:   end,
		},
		EnableMulti: enableMulti,
		Pi:          pi,
		TimeUse:     fmt.Sprintf("%vms", over.Sub(begin).Milliseconds()),
	}
	log(res)
}
