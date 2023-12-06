package server

// SimpleMessage contains a simple message for return.
var appID string = "cli_a5db390a6032500b"
var appSecret string = "67VN7qbvBG2kvdVcbosWqbld34UcvR7p"
var appToken string = `I2WUbQwTvaZCJIsQy3mcacxPndc`

// 定义 JSON 请求体的结构体
type TableData struct {
	Table struct {
		Name            string `json:"name"`
		DefaultViewName string `json:"default_view_name"`
		Fields          []struct {
			FieldName string `json:"field_name"`
			Type      int    `json:"type"`
		} `json:"fields"`
	} `json:"table"`
}

type ViewData struct {
	Name string `json:"view_name"`
	Type string `json:"view_type"`
}

type FieldData struct {
	Name string `json:"field_name"`
	Type int    `json:"type"`
}

type Attachment struct {
	Info map[string]string
}

// 生成图片的请求体
type ImageDescription struct {
	Description string `json:"description"`
}

// 新建记录的请求体
type MediaData struct {
	FileUrl string   `json:"FileUrl"`
	Prompt  string   `json:"Prompt"`
	Tags    []string `json:"Tags"`
}
