package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/big"
	"sort"
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
		var (
			number   *big.Int
			isNumber bool
		)

		resultTime = ""
		result = ""

		sNumber = ctx.PostForm("number")
		if number, isNumber = new(big.Int).SetString(sNumber, 10); !isNumber {
			result = "数字じゃない文字が混ざってるっぽいで"
		} else if number.Cmp(big.NewInt(0)) <= 0 {
			result = "0より大きい数字を入力してクレメンス"
		} else if number.Cmp(big.NewInt(1)) == 0 {
			result = "1は素数なんかな？ 少なくとも分解は出来んなあ。"
		} else {
			var results results
			z1Last.Set(big.NewInt(2))
			z2Last.Set(big.NewInt(2))

			if number.ProbablyPrime(20) {
				result = "この数字は素数やな"
				resultTime = ""
			} else {
				start := time.Now()
				for {
					calcFactor := calcFactor(number)
					results = append(results, calcFactor.String())
					number.Div(number, calcFactor)
					if number.ProbablyPrime(20) {
						results = append(results, number.String())
						break
					}
				}

				sort.Sort(results)

				end := time.Now()
				resultTime = "時間: " + fmt.Sprintf("%f秒\n", (end.Sub(start)).Seconds())
				result = "= " + strings.Join(results, " × ")
			}
		}

		ctx.Redirect(302, "/")
	})

	router.Run()
}

// calcFactor 素因数を計算
func calcFactor(n *big.Int) *big.Int {
	z1 := z1Last
	z2 := z2Last
	z2_z1 := big.NewInt(0)
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
		result.GCD(nil, nil, z2_z1.Sub(z2, z1).Abs(z2_z1), n)
		if result.Cmp(one) > 0 {
			return result
		}
	}

	for {
		z1.Mul(z1, z1).Add(z1, one).Mod(z1, n)
		z2.Mul(z2, z2).Add(z2, one).Mod(z2, n)
		z2.Mul(z2, z2).Add(z2, one).Mod(z2, n)

		result.GCD(nil, nil, z2_z1.Sub(z2, z1).Abs(z2_z1), n)
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

// isModZero 余りがゼロかどうかを判定 小さい素数用
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

// results 結果を文字列の配列で保存。sort用
type results []string

// Len 要素数。sort用
func (r results) Len() int {
	return len(r)
}

// Less iがjより小さくなる条件。sort用
func (r results) Less(i, j int) bool {
	iInt, _ := new(big.Int).SetString(r[i], 10)
	jInt, _ := new(big.Int).SetString(r[j], 10)

	return iInt.Cmp(jInt) < 0
}

// Swap 入れ替える方法。sort用
func (r results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
