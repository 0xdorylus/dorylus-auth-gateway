package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/carlmjohnson/requests"
	"github.com/form3tech-oss/jwt-go"
)

type AuthRequest struct {
	PublicKey string `json:"public_key"`
	Sig       string `json:"sig"`
	Ts        int64  `json:"ts"`
}

type AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type CheckNodeValidationRequest struct {
	PublicKey      string `json:"publicKey"`
	Signature      string `json:"signature"`
	SignContent    string `json:"signContent"`
	CheckPublicKey string `json:"checkPublicKey"`
}

type NodeInfo struct {
	PublicKey  string `json:"imei"`
	PrivateKey string `json:"sol_kp"`
	Ip         string `json:"ip"`
	CentralId  int    `json:"central_id"`
	KbType     int    `json:"kb_type"`
	account    types.Account
}

type CheckNodeVaildResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Check bool `json:"check"`
	} `json:"data"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var authReq AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusBadRequest,
			Msg:  "Bad request.",
		})
		io.WriteString(w, string(rsp))
		return
	}
	log.Printf("authReq: %v\n", authReq)
	now := time.Now().Unix()
	if authReq.Ts > now+MAX_DELTA_SECONDS || authReq.Ts < now-MAX_DELTA_SECONDS {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusUnauthorized,
			Msg:  "Time needs to be synchronized.",
		})
		io.WriteString(w, string(rsp))
		return
	}

	if !SolVerifySign(authReq.PublicKey, strconv.FormatInt(authReq.Ts, 10), authReq.Sig) {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusUnauthorized,
			Msg:  "Invalid signature.",
		})
		io.WriteString(w, string(rsp))
		return
	}

	// Get Node Info
	nodeInfo := GetNodeInfo()
	if nodeInfo == nil {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusInternalServerError,
			Msg:  "Node not prepared, please wait.",
		})
		io.WriteString(w, string(rsp))
		return
	}
	// log.Printf("nodeInfo: %v\n", nodeInfo)
	isNodeValid := nodeInfo.CheckNodeValidation(authReq.PublicKey)
	if !isNodeValid {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusUnauthorized,
			Msg:  "This node is reported invalid from scheduler server.",
		})
		io.WriteString(w, string(rsp))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"public_key": authReq.PublicKey,
		"exp":        time.Now().Add(time.Hour * time.Duration(JWT_TOKEN_EXPIRE_HOURS)).Unix(),
		"iat":        time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		rsp, _ := json.Marshal(&AuthResponse{
			Code: http.StatusInternalServerError,
			Msg:  "Token generation failed.",
		})
		io.WriteString(w, string(rsp))
		return
	}

	rsp, _ := json.Marshal(&AuthResponse{
		Code: 0,
		Data: tokenString,
	})
	io.WriteString(w, string(rsp))
}

func AuthMiddleware(next http.Handler) http.Handler {
	if len(APP_KEY) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(APP_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyValid := CheckApiKeyVaild(r)
		if apiKeyValid {
			next.ServeHTTP(w, r)
			return
		}

		err := jwtMiddleware.CheckJWT(w, r)
		if err != nil {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckApiKeyVaild(r *http.Request) bool {
	apiKeyHeader := r.Header.Get(API_KEY_HEADER)
	if apiKeyHeader == "" {
		return false
	}

	bindIp, ok := API_KEY_SET[apiKeyHeader]
	if !ok {
		return false
	}

	hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("failed to parse remote address: %v\n", r.RemoteAddr)
		return false
	}

	return bindIp == hostIP
}

var gNodeInfo *NodeInfo = nil
var loadNodeInfoOnce sync.Once

func GetNodeInfo() *NodeInfo {
	loadNodeInfoOnce.Do(LoadNodeInfo)
	return gNodeInfo
}

func LoadNodeInfo() {
	fileBytes, err := os.ReadFile(NODE_INFO_JSON_PATH)
	if err != nil {
		log.Printf("LoadNodeInfo from %v failed: %v\n", NODE_INFO_JSON_PATH, err)
		return
	}

	gNodeInfo = new(NodeInfo)
	err = json.Unmarshal(fileBytes, gNodeInfo)
	if err != nil {
		log.Printf("Unmarshal from %v failed: %v\n", NODE_INFO_JSON_PATH, err)
	}

	gNodeInfo.account, gNodeInfo.PublicKey, err = SolRestoreAccount(gNodeInfo.PrivateKey)
	if err != nil {
		log.Printf("SolRestoreAccount Error: %v\n", err)
	}
}

func (nodeInfo *NodeInfo) CheckNodeValidation(checkPublicKey string) bool {
	now := time.Now().UnixMilli()

	signContent := strconv.FormatInt(now, 10)
	sig := nodeInfo.account.Sign([]byte(signContent))
	sigStr := hex.EncodeToString(sig)

	body := &CheckNodeValidationRequest{
		PublicKey:      nodeInfo.PublicKey,
		Signature:      sigStr,
		SignContent:    signContent,
		CheckPublicKey: checkPublicKey,
	}

	var result CheckNodeVaildResult
	err := requests.URL(CHECK_NODE_VALID_URL).BodyJSON(&body).
		ToJSON(&result).Fetch(context.Background())
	if err != nil {
		log.Printf("CheckNodeVaild err: %v\n", err)
		return false
	}

	return result.Data.Check
}
