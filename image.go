package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	chatgpt_errors "github.com/gralexey/go-chatgpt/utils"
)

type ChatGPTImageModel string

const (
	Dalle2 ChatGPTImageModel = "dall-e-2"
	Dalle3 ChatGPTImageModel = "dall-e-3"
)

type ImageRequest struct {
	// (Required)
	// ID of the model to use.
	Model ChatGPTImageModel `json:"model"`

	// Required
	// The prompt to generate image for
	Prompt string `json:"prompt"`

	N int `json:"n"`

	Size ImageSize `json:"size"`
}

type ImageSize string

const (
	ImageSize1024x1024 ImageSize = "1024x1024"
	ImageSize512x512   ImageSize = "512x512"
)

type ImageResponse struct {
	Data []ImageResponseDataItem `json:"data"`
}

type ImageResponseDataItem struct {
	Url string `json:"url"`
}

func (c *Client) SimpleGenImage(ctx context.Context, prompt string) (*ImageResponse, error) {
	req := &ImageRequest{
		Model:  Dalle3,
		Prompt: prompt,
		Size:   ImageSize1024x1024,
		N:      1,
	}

	return c.SendImageGenRequest(ctx, req)
}

func (c *Client) SendImageGenRequest(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	if err := validateImageRequest(req); err != nil {
		return nil, err
	}

	reqBytes, _ := json.Marshal(req)

	endpoint := "/images/generations"
	httpReq, err := http.NewRequest("POST", c.config.BaseURL+endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	httpReq = httpReq.WithContext(ctx)

	res, err := c.sendRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var imageResponse ImageResponse
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&imageResponse); err != nil {
		return nil, err
	}

	return &imageResponse, nil
}

func validateImageRequest(req *ImageRequest) error {
	if len(req.Prompt) == 0 {
		return chatgpt_errors.ErrNoImagePrompt
	}

	isAllowed := false

	allowedModels := []ChatGPTImageModel{
		Dalle2, Dalle3,
	}

	for _, model := range allowedModels {
		if req.Model == model {
			isAllowed = true
		}
	}

	if !isAllowed {
		return chatgpt_errors.ErrInvalidModel
	}

	return nil
}
