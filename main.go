package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/big"
	"strings"
	"time"
)

// どこまで乱数が進んだかを記憶
var z1Last = big.NewInt(2)
var z2Last = big.NewInt(2)

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")
	result := ""
	var resultTime string
	var sNumber string

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{"number": sNumber, "result": result, "time": resultTime})
	})

	//Calc
	router.POST("/calc", func(ctx *gin.Context) {
		sNumber = ctx.PostForm("number")
		number, _ := new(big.Int).SetString(sNumber, 10)
		var arrayResults []string
		z1Last.Set(big.NewInt(2))
		z2Last.Set(big.NewInt(2))

		if number.ProbablyPrime(20) {
			result = "この数字は素数やな"
			resultTime = ""
		} else {
			start := time.Now()
			for {
				calcFactor := calcFactor(number)
				arrayResults = append(arrayResults, calcFactor.String())
				number.Div(number, calcFactor)
				if number.ProbablyPrime(20) {
					arrayResults = append(arrayResults, number.String())
					break
				}
			}

			end := time.Now()
			resultTime = "時間: " + fmt.Sprintf("%f秒\n", (end.Sub(start)).Seconds())
			result = "= " + strings.Join(arrayResults, " × ")
		}

		ctx.Redirect(302, "/")
	})

	router.Run()
}

func calcFactor(n *big.Int) *big.Int {
	z1 := z1Last
	z2 := z2Last
	z2_z1 := big.NewInt(0)
	a := big.NewInt(1)
	b := big.NewInt(1)
	result := big.NewInt(1)
	one := big.NewInt(1)
	numbers := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(7), big.NewInt(11), big.NewInt(13)}

	// 小さい値が素因数だとたまにエラーが起こるので対策
	for _, number := range numbers {
		if isModZero(n, number) {
			return number
		}
	}

	if z1.Cmp(z2) != 0 {
		result.GCD(a, b, z2_z1.Sub(z2, z1).Abs(z2_z1), n)
		if result.Cmp(one) > 0 {
			return result
		}
	}

	for {
		z1.Mul(z1, z1)
		z1.Add(z1, one)
		z1.Mod(z1, n)
		z2.Mul(z2, z2)
		z2.Add(z2, one)
		z2.Mod(z2, n)
		z2.Mul(z2, z2)
		z2.Add(z2, one)
		z2.Mod(z2, n)

		result.GCD(a, b, z2_z1.Sub(z2, z1).Abs(z2_z1), n)
		if result.Cmp(one) > 0 {
			break
		}
	}
	if result.ProbablyPrime(20) {
		z1Last.Set(z1)
		z2Last.Set(z2)
		return result
	} else {
		return calcFactor(result)
	}
}

// isModZero 余りがゼロかどうかを判定
func isModZero(n, m *big.Int) bool {
	zero := big.NewInt(0)
	module := big.NewInt(0)
	n.DivMod(n, m, module)
	if module.Cmp(zero) == 0 {
		n.Mul(n, m)
		return true
	}
	n.Mul(n, m).Add(n, module)
	return false
}
