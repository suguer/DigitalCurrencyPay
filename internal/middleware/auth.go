package middleware

import (
	"DigitalCurrency/internal/service/user"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Html404 = []byte(`<html>
<head><title>404 Not Found</title></head>
<body bgcolor="white">
<center><h1>404 Not Found</h1></center>
<hr><center>AllinSSL</center>
</body>
</html>`)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		if c.Request.Method != "POST" {
			c.AbortWithStatus(405)
			return
		}

		// 获取请求中的所有参数
		var requestData map[string]any
		err := c.ShouldBindBodyWith(&requestData, binding.JSON)
		if err != nil {
			c.ShouldBindQuery(&requestData)
		}

		// 获取签名
		sign, exists := requestData["sign"]
		if !exists {
			c.JSON(400, gin.H{"error": "缺少签名参数"})
			c.Abort()
			return
		}
		merchant_id, exists := requestData["merchant_id"]
		if !exists {
			c.JSON(400, gin.H{"error": "缺少商户ID参数"})
			c.Abort()
			return
		}
		user, err := user.GetSecret(merchant_id)
		if err != nil {
			c.JSON(401, gin.H{"error": "商户ID不存在"})
			c.Abort()
			return
		}

		secret := user.Secret // 实际应用中应该从配置中获取

		// 计算签名
		calculatedSign := Signature(requestData, secret)
		// 验证签名
		if calculatedSign != sign {
			c.JSON(401, gin.H{"error": "签名验证失败"})
			c.Abort()
			return
		}

		// 签名验证通过，设置用户ID并继续
		var user_id uint = user.ID // 实际应用中应该根据验证结果获取用户ID
		c.Set("user_id", user_id)
	}
}

// Signature generates a HMAC SHA256 signature for the given data and secret
// This is a Golang implementation of the PHP signature function
func Signature(data map[string]any, secret string) string {
	// Create a slice of keys for sorting
	keys := make([]string, 0, len(data))
	for k := range data {
		if k != "sign" {
			keys = append(keys, k)
		}
	}

	// Sort the keys
	sort.Strings(keys)

	// Build the signature string
	var signBuilder strings.Builder
	for _, k := range keys {
		signBuilder.WriteString(k)
		signBuilder.WriteString("=")
		// Convert any value type to string
		switch v := data[k].(type) {
		case string:
			signBuilder.WriteString(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			// 整数值使用%d格式化，避免科学计数法
			signBuilder.WriteString(fmt.Sprintf("%d", v))
		case float32, float64:
			// 如果是浮点数，使用%f格式化并去除尾随零
			s := fmt.Sprintf("%f", v)
			s = strings.TrimRight(s, "0")
			s = strings.TrimRight(s, ".")
			signBuilder.WriteString(s)
		default:
			// 对于其他类型，尝试转换为数字处理
			if num, ok := v.(json.Number); ok {
				if i, err := num.Int64(); err == nil {
					signBuilder.WriteString(fmt.Sprintf("%d", i))
				} else if f, err := num.Float64(); err == nil {
					if f == float64(int64(f)) {
						signBuilder.WriteString(fmt.Sprintf("%d", int64(f)))
					} else {
						s := fmt.Sprintf("%f", f)
						s = strings.TrimRight(s, "0")
						s = strings.TrimRight(s, ".")
						signBuilder.WriteString(s)
					}
				} else {
					signBuilder.WriteString(num.String())
				}
			} else {
				signBuilder.WriteString(fmt.Sprintf("%v", v))
			}
		}
		signBuilder.WriteString("&")
	}
	signBuilder.WriteString("secret=")
	signBuilder.WriteString(secret)

	signString := signBuilder.String()

	// Create HMAC SHA256 hash
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signString))

	// Return the hex-encoded hash
	return hex.EncodeToString(h.Sum(nil))
}
