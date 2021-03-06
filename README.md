# jwc
[DEPRECATED]
I'd recommend https://github.com/square/go-jose instead

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/khezen/jwc)
[![Build Status](http://img.shields.io/travis/khezen/jwc.svg?style=flat-square)](https://travis-ci.org/khezen/jwc) [![codecov](https://img.shields.io/codecov/c/github/khezen/jwc/master.svg?style=flat-square)](https://codecov.io/gh/khezen/jwc)
[![Go Report Card](https://goreportcard.com/badge/github.com/khezen/jwc?style=flat-square)](https://goreportcard.com/report/github.com/khezen/jwc)

JSON Web Cryptography

* jws - JSON Web Signature
  * RSASSA-PSS + SHA256, recommended +
  * RSASSA-PKCS1-v1_5 + SHA256, recommended -
  
* jwe - JSON Web Encryption
  * key encryption
    * RSA-OAEP, recommended +
    * RSAES-PKCS1-v1_5, recommended -
  * content encryption
    * A128CBC-HS256
    * A192CBC-HS384
    * A256CBC-HS512
    * A128GCM
    * A192GCM
    * A256GCM


## Example

### Signature

```golang
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/khezen/jwc"
)

func main() {
	var (
		privateKey, _ = rsa.GenerateKey(rand.Reader, 1024)
		jwkid         = jwc.JWKID("52d510e3-8d0a-4ef8-a81a-c8cd7ce06472")
		jwk, _        = jwc.RSAToPublicJWK(&privateKey.PublicKey, jwkid, jwc.PS256, nil)
		jwkBytes, _   = json.Marshal(jwk)
	)
	jwtStr := issueJWT(jwkid, privateKey)
	token := verify(jwtStr, jwkBytes)
	fmt.Println(jwtStr)
	fmt.Println()
	fmt.Println(token)
}

func issueJWT(keyID jwc.JWKID, privateKey *rsa.PrivateKey) string {
	now := time.Now().UTC()
	nowUnix := now.Unix()
	exp := now.Add(time.Minute)
	expUnix := exp.Unix()
	jwt, err := jwc.NewJWT(
		jwc.JWTPayload{
			RegisteredClaims: jwc.RegisteredClaims{
				IssuedAtTimestamp:   nowUnix,
				ExpirationTimestamp: expUnix,
				Issuer:              "github.com/khezen/jwc/jwt_test.go",
				Subject:             "customer_id",
				Audiance:            "android.myapp.com",
			},
			PrivateClaims: jwc.PrivateClaims{
				"id":  "token_id",
				"did": "device_id",
				"sco": "offline",
				"cc":  "dummy_code_challenge",
				"ccm": "S256",
			},
		},
		jwc.PS256,
	)
	if err != nil {
		panic(err)
	}
	jwtStr, err := jwt.Encode(keyID, privateKey)
	if err != nil {
		panic(err)
	}
	return jwtStr
}

func verify(jwtStr string, jwkBytes []byte) *jwc.JWT {
	var pubJWK jwc.RSAPublicJWK
	err := json.Unmarshal(jwkBytes, &pubJWK)
	if err != nil {
		panic(err)
	}
	pubKey, err := pubJWK.PublicRSA()
	if err != nil {
		panic(err)
	}
	token, err := jwc.DecodeVerifyJWT(jwtStr, pubKey)
	if err != nil {
		panic(err)
	}
	return token
}
```

#### output

```sh
eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjRkYjY2NzgwLWYzNjMtNDdmZC05MDlkLTQ0MWM1ZmUyM2Q2MSJ9.eyJleHAiOjE1NDYzNjk0NTIsImlhdCI6MTU0NjM2OTM5MiwiaXNzIjoiZ2l0aHViLmNvbS9raGV6ZW4vandjL2p3dF90ZXN0LmdvIiwic3ViIjoidGVzdCIsInByaXYiOnsiYXVkIjoiYW5kcm9pZC5teWFwcC5jb20iLCJjYyI6ImR1bW15Q29kZUNoYWxsZW5nZSIsImNjbSI6IlMyNTYiLCJjaWQiOiIwMTA4YmZiZC0yZjc2LTQyMmEtYTNiZC0yNjMxZmVhYWNiZWUiLCJkaWQiOiJkZXZpY2VJRCIsInNjbyI6Im9mZmxpbmUiLCJ0aWQiOiJiMGIwZWM5Mi1jZTNjLTQ3ZjUtODQ5Ny03Y2FiMjkxNDcyZDAifX0.WktA5tt_Tt6R-qZuTqpSB7xnYDrMlJXjz7aTzQys1UjMAEjLFHCWqmLp33DRlUboZiZQWa_6D4c6fzS-UHFQ9pQ_73s_Rg83i6XEMJIlr2k420g_cO-N_y425gnoJ2GDOpVSGxMS5uofh8JoE6OZpPNauJo_Z5MNpEKp5XZDEAE

&{{PS256 JWT 4db66780-f363-47fd-909d-441c5fe23d61} {{1546369452 1546369392 github.com/khezen/jwc/jwt_test.go test } map[tid:b0b0ec92-ce3c-47f5-8497-7cab291472d0 aud:android.myapp.com cc:dummyCodeChallenge ccm:S256 cid:0108bfbd-2f76-422a-a3bd-2631feaacbee did:deviceID sco:offline]}}
```

### Encryption

```golang
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/khezen/jwc"
)

func main() {
	var (
		privateKey, _     = rsa.GenerateKey(rand.Reader, 2042)
		jwkid         = jwc.JWKID("52d510e3-8d0a-4ef8-a81a-c8cd7ce06472")
		publicJWK, _      = jwc.RSAToPublicJWK(&privateKey.PublicKey, jwkid, jwc.ROAEP, nil)
		publicJWKBytes, _ = json.Marshal(publicJWK)
		message           = []byte("lorem ipsum ipsa occaecati aut velit facilis enim dolorum id eius magni ducimus sed illum similique cupiditate sit id perferendis alias sint")
	)
	compactJWE := encrypt(message, publicJWKBytes)
	plaintext := decrypt(compactJWE, privateKey)
	fmt.Println(bytes.EqualFold(message, plaintext))
}

func encrypt(plaintext, jwkBytes []byte) []byte {
	var jwk jwc.RSAPublicJWK
	err := json.Unmarshal(jwkBytes, &jwk)
	if err != nil {
		panic(err)
	}
	pubKey, err := jwk.PublicRSA()
	if err != nil {
		panic(err)
	}
	jwe, err := jwc.NewJWE(
		&jwc.JOSEHeaders{Algorithm: jwc.ROAEP, Encryption: jwc.A256GCM},
		pubKey,
		plaintext,
	)
	if err != nil {
		panic(err)
	}
	jweString, err := jwe.Compact()
	if err != nil {
		panic(err)
	}
	return jweString
}

func decrypt(compactJWE []byte, privateKey *rsa.PrivateKey) []byte {
	jwe, err := jwc.ParseCompactJWE(compactJWE)
	if err != nil {
		panic(err)
	}
	plaintext, err := jwe.Plaintext(privateKey)
	if err != nil {
		panic(err)
	}
	return plaintext
}
```

#### output

```sh
true
```