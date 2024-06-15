package main

const (
	APP_KEY                = "NodeAuth@dorylus.ai"
	MAX_DELTA_SECONDS      = 30
	NODE_INFO_JSON_PATH    = "/opt/dorylus/cache/node_info.json"
	CHECK_NODE_VALID_URL   = "https://hub.dorylus.ai/api/restful/api/checkNode"
	JWT_TOKEN_EXPIRE_HOURS = 1
	API_KEY_HEADER         = "X-API-Key"
)

var API_KEY_SET = map[string]string{
	"hjm-059da3a859c48f563671191c34d1fe87d558622d": "10.10.0.1",
	"hjm-dbc063418fc09af4d81046df69a28f4a3e94fdb8": "10.10.0.2",
}
