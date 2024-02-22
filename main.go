package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"

    "github.com/go-resty/resty/v2"
)

const (
    apiEndpoint = "https://api.openai.com/v1/chat/completions"
    configPath  = "config.json"
)

type Config struct {
    APIKey string `json:"apiKey"`
}

func loadConfig(path string) (Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return Config{}, err
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return Config{}, err
    }

    return config, nil
}


func main() {
    // Read API key from config file
    config, err := loadConfig(configPath)
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    client := resty.New()

    response, err := client.R().
        SetAuthToken(config.APIKey).
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "model":      "gpt-3.5-turbo",
            "messages":   []interface{}{map[string]interface{}{"role": "system", "content": "Hi can you tell me what is the factorial of 10?"}},
            "max_tokens": 50,
        }).
        Post(apiEndpoint)

    if err != nil {
        log.Fatalf("Error while sending the request: %v", err)
    }

    body := response.Body()

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fmt.Println("Error while decoding JSON response:", err)
        return
    }

	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
	    fmt.Println("No choices found in the response")
	    return
	}

	// Now you can safely access the first choice
	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	fmt.Println(content)
}

