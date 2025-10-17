- 客户端在调用API时,需要在请求头中添加签名信息
- 签名信息包含以下字段:
    - `merchant_id`: 时间戳,单位秒
    - `timestamp`: 时间戳,单位秒
    - `nonce`: 随机数,用于防止重放攻击
    - `sign_type`: 签名类型,默认值为`HMAC-SHA256`
    - `sign`: 签名,使用HMAC-SHA256算法计算
- 签名计算方式:
    - 将请求参数按照字典序排序
    - 忽略`sign`参数
    - 拼接参数为`key1=value1&key2=value2&...&keyN=valueN`格式
    - 最后在拼接字符串末尾添加`secret_key`参数,格式为`secret=secret_key`
    - 对拼接后的字符串进行HMAC-SHA256加密
    - 最后将加密后的结果转换为十六进制字符串


### php实例
```php 

   public function request($action, $params = [])
    {
        $origin = $params;
        $params['merchant_id'] = $this->merchantId;
        $params['sign_type'] = 'HMAC-SHA256';
        $params['timestamp'] = time();
        $params['sign'] = $this->signature($params, $this->secret);

        try {
            $http = new Client(['timeout' => 30]);
            $response = $http->post($action, ['json' => $params]);
            $body = $response->getBody()->getContents();
            $json = json_decode($body, true);
        } catch (\Exception $e) {
            $json = ['code' => 500, 'message' => $e->getMessage()];
        }
        return $json;
    }


    public function signature($data, $secret)
    {
        ksort($data);

        $signString = '';
        foreach ($data as $key => $value) {
            if ($key == 'sign') continue;

            $signString .= $key . '=' . $value . '&';
        }
        $signString .= 'secret=' . $secret;

        return hash_hmac('sha256', $signString, $secret);
    }
```

### golang实例


```golang

func Signature(data map[string]interface{}, secret string) string {
	// 按键名排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	signString := ""
	for _, key := range keys {
		if key == "sign" {
			continue
		}
		
		// 值转换为字符串
		var strValue string
		switch v := data[key].(type) {
		case string:
			strValue = v
		case int:
			strValue = strconv.Itoa(v)
		case int64:
			strValue = strconv.FormatInt(v, 10)
		default:
			strValue = toString(v)
		}
		
		signString += key + "=" + strValue + "&"
	}
	
	signString += "secret=" + secret

	// 计算 HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signString))
	return hex.EncodeToString(h.Sum(nil))
}

// Request 准备请求参数
func Request(merchantID, secret, action string, params map[string]interface{}) map[string]interface{} {
	// 设置固定参数
	params["merchant_id"] = merchantID
	params["sign_type"] = "HMAC-SHA256"
	params["timestamp"] = time.Now().Unix()
	
	// 生成签名
	params["sign"] = Signature(params, secret)
	
	return params
}

// toString 通用类型转字符串
func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}
```