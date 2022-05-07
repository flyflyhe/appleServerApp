package apple

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/dvsekhvalnov/jose2go/compact"
	"github.com/dvsekhvalnov/jose2go/keys/ecc"
	"github.com/flyflyhe/appleServerApp/config"
	"github.com/flyflyhe/appleServerApp/services/arrayHelper"
	"github.com/syyongx/php2go"
	"gopkg.in/yaml.v2"
)

type LookUp struct {
	Status             int      `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

type TransactionHistory struct {
	Revision           string   `json:"revision"`
	BundleId           string   `json:"bundleId"`
	AppAppleId         int      `json:"appAppleId"`
	Environment        string   `json:"environment"`
	HasMore            bool     `json:"hasMore"`
	SignedTransactions []string `json:"signedTransactions"`
}

type AppleHeaders struct {
	Alg string   `json:"alg"`
	X5c []string `json:"x5c"`
}

type AppleJwtConfig struct {
	Kid string `yaml:"kid"`
	Iss string `yaml:"iss"`
	Bid string `yaml:"bid"`
}

var privateKey string

var appleKey string

var appleJwtConfig AppleJwtConfig

func init() {

	// 转换成Struct
	if err := yaml.Unmarshal(config.GetYamlConfig(), &appleJwtConfig); err != nil {
		if err != nil {
			panic(err)
		}
	}

	privateKey = string(config.GetApplePrivateKey())
	appleKey = string(config.GetAppleRootCaPem())
}

type ErrMsg struct {
	Err string
}

func NewErrMsg(str string) *ErrMsg {
	return &ErrMsg{Err: str}
}

func (errMsg ErrMsg) Error() string {
	return errMsg.Err
}

func CheckOrder(orderId, token string) (result []string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = NewErrMsg(err2.(string)) //这里
		}
	}()

	url := "https://api.storekit.itunes.apple.com/inApps/v1/lookup/" + orderId

	var request *http.Request
	if request, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		panic(err)
	}

	request.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	lookUp := &LookUp{}
	if err := json.Unmarshal(body, lookUp); err != nil {
		panic(err)
	}

	if lookUp.Status == 1 {
		log.Printf("非法订单号")
		return nil, NewErrMsg("非法订单号")
	}

	for _, transaction := range lookUp.SignedTransactions {
		parts, _ := compact.Parse(transaction)

		headers := parts[0]
		//payload, _ := base64url.Decode(string(parts[1]))
		//signature, _ := base64url.Decode(string(parts[2]))

		appleHeader := &AppleHeaders{}

		if err := json.Unmarshal(headers, appleHeader); err != nil {
			panic(err)
		}

		//证书链校验 c1 -> c2 -> c3 -> apple
		certificates := make([]string, 0)
		for _, v := range appleHeader.X5c {
			certificate := "-----BEGIN CERTIFICATE-----\n"
			certificate += php2go.ChunkSplit(v, 64, "\n")
			certificate += "-----END CERTIFICATE-----\n"
			certificates = append(certificates, certificate)
		}

		certificates = append(certificates, appleKey)

		var certFirst *x509.Certificate

		for i := 0; i < len(certificates); i++ {
			block, _ := pem.Decode([]byte(certificates[i]))
			if block == nil {
				panic("failed to parse certificate PEM")
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				panic("failed to parse certificate: " + err.Error())
			}

			if certFirst == nil {
				certFirst = cert
				continue
			}

			if err = certFirst.CheckSignatureFrom(cert); err != nil {
				panic(err)
			}

			certFirst = cert
		}

		//decode 传入第一个证书
		var ecdsaPublic *ecdsa.PublicKey
		if ecdsaPublic, err = ecc.ReadPublic([]byte(certificates[0])); err != nil {
			panic(err)
		}

		t, _, err := jose.Decode(transaction, ecdsaPublic)
		if err != nil {
			panic(err)
		}

		result = append(result, t)
	}

	return result, nil
}

func GetTransactionHistory(originalTransactionId, token, revision string) (result []string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			fmt.Println(err2)
			err = NewErrMsg("异常 订单号有误或重试") //这里
		}
	}()

	url := "https://api.storekit.itunes.apple.com/inApps/v1/history/" + originalTransactionId

	if revision != "" {
		url += "?revision=" + revision
	}

	var request *http.Request
	if request, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		panic(err)
	}

	request.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(res.Body)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	transactionHistory := &TransactionHistory{}
	if err := json.Unmarshal(body, transactionHistory); err != nil {
		panic(err)
	}

	if transactionHistory.HasMore {
		resultTmp, err := GetTransactionHistory(originalTransactionId, token, transactionHistory.Revision)
		if err != nil {
			panic(err)
		}
		arrayHelper.ArrayReverse(resultTmp)
		result = append(result, resultTmp...)
	}

	for _, transaction := range transactionHistory.SignedTransactions {
		parts, _ := compact.Parse(transaction)

		headers := parts[0]
		//payload, _ := base64url.Decode(string(parts[1]))
		//signature, _ := base64url.Decode(string(parts[2]))

		appleHeader := &AppleHeaders{}

		if err := json.Unmarshal(headers, appleHeader); err != nil {
			panic(err)
		}

		//证书链校验 c1 -> c2 -> c3 -> apple
		certificates := make([]string, 0)
		for _, v := range appleHeader.X5c {
			certificate := "-----BEGIN CERTIFICATE-----\n"
			certificate += php2go.ChunkSplit(v, 64, "\n")
			certificate += "-----END CERTIFICATE-----\n"
			certificates = append(certificates, certificate)
		}

		certificates = append(certificates, appleKey)

		var certFirst *x509.Certificate

		for i := 0; i < len(certificates); i++ {
			block, _ := pem.Decode([]byte(certificates[i]))
			if block == nil {
				panic("failed to parse certificate PEM")
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				panic("failed to parse certificate: " + err.Error())
			}

			if certFirst == nil {
				certFirst = cert
				continue
			}

			if err = certFirst.CheckSignatureFrom(cert); err != nil {
				panic(err)
			}

			certFirst = cert
		}

		//decode 传入第一个证书
		var ecdsaPublic *ecdsa.PublicKey
		if ecdsaPublic, err = ecc.ReadPublic([]byte(certificates[0])); err != nil {
			panic(err)
		}

		t, _, err := jose.Decode(transaction, ecdsaPublic)
		if err != nil {
			panic(err)
		}

		result = append(result, t)
	}

	return result, nil
}

func GetAppleJwtToken() string {
	headerConfig := map[string]interface{}{
		"alg": "ES256",
		"kid": appleJwtConfig.Kid,
		"typ": "JWT",
	}

	joseConfig := jose.Headers(headerConfig)

	now := time.Now()
	payloadConfig := map[string]interface{}{
		"iss": appleJwtConfig.Iss,
		"iat": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),
		"aud": "appstoreconnect-v1",
		"bid": appleJwtConfig.Bid,
	}

	payload, _ := json.Marshal(payloadConfig)

	var token string
	var err error

	var ecdsaPrivate *ecdsa.PrivateKey
	if ecdsaPrivate, err = ecc.ReadPrivate([]byte(privateKey)); err != nil {
		panic(err)
	}

	if token, err = jose.SignBytes(payload, jose.ES256, ecdsaPrivate, joseConfig); err != nil {
		panic(err)
	}

	return token
}
