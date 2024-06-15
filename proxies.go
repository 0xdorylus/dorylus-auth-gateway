package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pretty66/websocketproxy"
)

func OllamaLLMHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/api/generate"
	u, _ := url.Parse("http://127.0.0.1:11434/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func OllamaLLMChatHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/api/chat"
	u, _ := url.Parse("http://127.0.0.1:11434/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func KbSearchHandler(w http.ResponseWriter, r *http.Request) {
	publicKey := strings.TrimPrefix(r.URL.Path, "/kb_search/")
	r.URL.Path = fmt.Sprintf("/api/search_docs/%v", publicKey)
	u, _ := url.Parse("http://127.0.0.1:5101/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func TextToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/tts"
	// body, err := io.ReadAll(r.Body)
	// r.Body = io.NopCloser(bytes.NewBuffer(body))

	u, _ := url.Parse("http://127.0.0.1:5000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func SDGetModelsHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/sdapi/v1/sd-models"
	u, _ := url.Parse("http://127.0.0.1:7860/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func SDGetSamplersHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/sdapi/v1/samplers"
	u, _ := url.Parse("http://127.0.0.1:7860/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func SDGetLorasHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/sdapi/v1/loras"
	u, _ := url.Parse("http://127.0.0.1:7860/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func SDTextToImageHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/sdapi/v1/txt2img"
	u, _ := url.Parse("http://127.0.0.1:7860/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func BlipGenerateHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/blip/api/generate"
	u, _ := url.Parse("http://127.0.0.1:7801/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func ClipClassifierHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/clip/api/classifier"
	u, _ := url.Parse("http://127.0.0.1:7802/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func ComparingCaptioningHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/comparing_captioning/api/generate"
	u, _ := url.Parse("http://127.0.0.1:7803/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func WsAsrStreamHandler() func(w http.ResponseWriter, r *http.Request) {
	wp, err := websocketproxy.NewProxy("ws://127.0.0.1:5102/asr/stream", func(r *http.Request) error {
		return nil
	})
	if err != nil {
		return nil
	}

	return wp.Proxy
}

func DorylusAgentClear(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/api/dorylus_agent_clear"
	u, _ := url.Parse("http://127.0.0.1:5101/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func DorylusAgentText(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/api/dorylus_agent_text"
	u, _ := url.Parse("http://127.0.0.1:5101/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func DorylusAgentAudio(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/api/dorylus_agent_audio"
	u, _ := url.Parse("http://127.0.0.1:5101/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
