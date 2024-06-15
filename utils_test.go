package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/carlmjohnson/requests"
)

func TestSendAuthRequest(t *testing.T) {
	privateKey := "4NGo9CkzcbingicrSTjaHSL3hgV2xLXMVWxFXQfMRLz2zy31QQa33HK3cmLkfUF7G4wLQKZCAz8EUVuHT5SK7HYN"
	account, err := types.AccountFromBase58(privateKey)
	if err != nil {
		t.Fatalf("Invaild privateKey: %v\n", privateKey)
	}
	publicKey := account.PublicKey.ToBase58()
	t.Logf("publicKey: %v\n", publicKey)

	now := time.Now().Unix()
	t.Logf("now: %v\n", now)

	sig := account.Sign([]byte(strconv.FormatInt(now, 10)))
	sigStr := hex.EncodeToString(sig)
	t.Logf("sig: %v\n", sigStr)

	url := "http://127.0.0.1:8080/auth"
	body := &AuthRequest{
		PublicKey: publicKey,
		Sig:       sigStr,
		Ts:        now,
	}
	t.Logf("body: %v\n", body)

	authRsp := new(AuthResponse)
	err = requests.URL(url).BodyJSON(&body).
		ToJSON(&authRsp).Fetch(context.Background())
	if err != nil {
		t.Fatalf("err: %v\n", err)
	} else {
		jsonBytes, _ := json.Marshal(authRsp)
		t.Logf("authRsp: %v\n", string(jsonBytes))
	}

	t.Log("TestSendAuthRequest ok")
}

func TestSendLLMRequest(t *testing.T) {

	body := &struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  "llama3",
		Prompt: "Hello.",
		Stream: false,
	}
	url := "http://127.0.0.1:8080/llm"
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU5MTUwOTQsImlhdCI6MTcxNTkxMTQ5NCwicHVibGljX2tleSI6IkdISFNXbVk1VG04TUc3Z2tNNXh6VjNrNzdvUVBtMVQyb1FaVjVlMTVYRHZpIn0.TuqjhmI8tPKGCRvQavO_QHqzCoyAHRrmfqKts0suku0"

	var rsp string
	err := requests.URL(url).BodyJSON(&body).Bearer(jwtToken).
		ToString(&rsp).Fetch(context.Background())
	if err != nil {
		t.Fatalf("err: %v\n", err)
	} else {
		t.Logf("rsp: %v\n", rsp)
	}

	t.Log("TestSendLLMRequest ok")
}

func TestSendKbSearchRequest(t *testing.T) {

	body := &struct {
		KbName string `json:"kb_name"`
		Query  string `json:"query"`
	}{
		KbName: "kb_1",
		Query:  "memory limits",
	}

	url := "http://127.0.0.1:8080"
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU4Njg3NTIsImlhdCI6MTcxNTg2NTE1MiwicHVibGljX2tleSI6IkdISFNXbVk1VG04TUc3Z2tNNXh6VjNrNzdvUVBtMVQyb1FaVjVlMTVYRHZpIn0.4rnH0P_2zqfifHtQ4cK4QgQ5iqhIltCbK_gJuZhwnag"
	publicKey := "GHHSWmY5Tm8MG7gkM5xzV3k77oQPm1T2oQZV5e15XDvi"

	var rsp string
	err := requests.URL(url).Pathf("/kb_search/%v", publicKey).BodyJSON(&body).Bearer(jwtToken).
		ToString(&rsp).Fetch(context.Background())
	if err != nil {
		t.Fatalf("err: %v\n", err)
	} else {
		t.Logf("rsp: %v\n", rsp)
	}

	t.Log("TestSendLLMRequest ok")
}

func TestGatewaySendLLMRequest(t *testing.T) {

	body := &struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  "llama3",
		Prompt: "Hello.",
		Stream: false,
	}
	url := "http://127.0.0.1:8080/llm"
	var apiKey string
	for k, _ := range API_KEY_SET {
		apiKey = k
		break
	}

	var rsp string
	err := requests.URL(url).BodyJSON(&body).Header(API_KEY_HEADER, apiKey).
		ToString(&rsp).Fetch(context.Background())
	if err != nil {
		t.Fatalf("err: %v\n", err)
	} else {
		t.Logf("rsp: %v\n", rsp)
	}

	t.Log("TestGatewaySendLLMRequest ok")
}

func TestGatewaySendKbSearchRequest(t *testing.T) {

	body := &struct {
		KbName string `json:"kb_name"`
		Query  string `json:"query"`
	}{
		KbName: "kb_1",
		Query:  "memory limits",
	}

	url := "http://127.0.0.1:8080"
	publicKey := "GHHSWmY5Tm8MG7gkM5xzV3k77oQPm1T2oQZV5e15XDvi"
	var apiKey string
	for k, _ := range API_KEY_SET {
		apiKey = k
		break
	}

	var rsp string
	err := requests.URL(url).Pathf("/kb_search/%v", publicKey).BodyJSON(&body).Header(API_KEY_HEADER, apiKey).
		ToString(&rsp).Fetch(context.Background())
	if err != nil {
		t.Fatalf("err: %v\n", err)
	} else {
		t.Logf("rsp: %v\n", rsp)
	}

	t.Log("TestGatewaySendKbSearchRequest ok")
}

func TestCheckNodeValidation(t *testing.T) {
	checkPublicKey := "GHHSWmY5Tm8MG7gkM5xzV3k77oQPm1T2oQZV5e15XDvi"
	privateKey := "LRSuhrSbg9yEVxHsrK2hYo1Wibw3bWM6czYDdv4pt81TPQTMUJK9ywpAkRG6pRfC2hNVf1CZ53k1aa1ZRtVNXCW"
	nodeInfo := &NodeInfo{
		PrivateKey: privateKey,
	}
	var err error
	nodeInfo.account, nodeInfo.PublicKey, err = SolRestoreAccount(nodeInfo.PrivateKey)
	if err != nil {
		t.Fatalf("SolRestoreAccount Error: %v\n", err)
	}
	result := nodeInfo.CheckNodeValidation(checkPublicKey)
	t.Logf("TestCheckNodeValidation result: %v\n", result)
}

func TestLoadNodeInfo(t *testing.T) {
	nodeInfo := GetNodeInfo()
	t.Logf("nodeInfo: %v\n", nodeInfo)
}
