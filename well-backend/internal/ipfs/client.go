package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Client struct {
	apiURL string
	hc     *http.Client
}

func New(apiURL string) *Client {
	return &Client{
		apiURL: strings.TrimRight(apiURL, "/"),
		hc:     &http.Client{},
	}
}

type addResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func (c *Client) Add(ctx context.Context, data []byte) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", "record.bin")
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := fw.Write(data); err != nil {
		return "", fmt.Errorf("write form: %w", err)
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+"/api/v0/add?pin=true", &b)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("ipfs add: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ipfs add status=%d body=%s", resp.StatusCode, string(body))
	}

	var result addResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ipfs add decode: %w", err)
	}

	return result.Hash, nil
}

func (c *Client) Cat(ctx context.Context, cid string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v0/cat?arg=%s", c.apiURL, cid)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ipfs cat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ipfs cat status=%d body=%s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ipfs cat read: %w", err)
	}

	return data, nil
}

func (c *Client) Pin(ctx context.Context, cid string) error {
	url := fmt.Sprintf("%s/api/v0/pin/add?arg=%s", c.apiURL, cid)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("ipfs pin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ipfs pin status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) Unpin(ctx context.Context, cid string) error {
	url := fmt.Sprintf("%s/api/v0/pin/rm?arg=%s", c.apiURL, cid)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("ipfs unpin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ipfs unpin status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}
