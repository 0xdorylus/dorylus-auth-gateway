package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/auth", AuthHandler)
	http.Handle("/llm/api/generate", AuthMiddleware(http.HandlerFunc(OllamaLLMHandler)))
	http.Handle("/llm/api/chat", AuthMiddleware(http.HandlerFunc(OllamaLLMChatHandler)))
	http.Handle("/kb_search/", AuthMiddleware(http.HandlerFunc(KbSearchHandler)))
	http.Handle("/tts", AuthMiddleware(http.HandlerFunc(TextToSpeechHandler)))

	http.Handle("/sd/models", AuthMiddleware(http.HandlerFunc(SDGetModelsHandler)))
	http.Handle("/sd/samplers", AuthMiddleware(http.HandlerFunc(SDGetSamplersHandler)))
	http.Handle("/sd/loras", AuthMiddleware(http.HandlerFunc(SDGetLorasHandler)))
	http.Handle("/sd/txt2img", AuthMiddleware(http.HandlerFunc(SDTextToImageHandler)))

	http.Handle("/blip/api/generate", AuthMiddleware(http.HandlerFunc(BlipGenerateHandler)))
	http.Handle("/clip/api/classifier", AuthMiddleware(http.HandlerFunc(ClipClassifierHandler)))
	http.Handle("/comparing_captioning/api/generate", AuthMiddleware(http.HandlerFunc(ComparingCaptioningHandler)))

	http.Handle("/asr/stream", AuthMiddleware(http.HandlerFunc(WsAsrStreamHandler())))

	http.Handle("/agent/api/dorylus_agent_clear", AuthMiddleware(http.HandlerFunc(DorylusAgentClear)))
	http.Handle("/agent/api/dorylus_agent_text", AuthMiddleware(http.HandlerFunc(DorylusAgentText)))
	http.Handle("/agent/api/dorylus_agent_audio", AuthMiddleware(http.HandlerFunc(DorylusAgentAudio)))

	// Start a basic HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
