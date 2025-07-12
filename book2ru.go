package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	version     = "0.1.0-go"
	httpReferer = "https://github.com/book2ru-go"
	xTitle      = "book2ru-go CLI (Go)"
	userAgent   = "book2ru-go-go/0.1.0"
)

var (
	batchSize = 10000 // Batch size in bytes
	apiURL    = "https://openrouter.ai/api/v1/chat/completions"
)

// Config holds the application configuration.
type Config struct {
	Model          string `yaml:"model"`
	Prompt         string `yaml:"prompt"`
	MetadataFooter bool   `yaml:"metadata_footer"`
	RetryAttempts  int    `yaml:"retry_attempts"`
	StartBatch     int    `yaml:"start_batch"`
	APIKey         string
}

// Batch represents a batch of text to translate.
type Batch struct {
	Content string
	Size    int
	Lines   int
}

// APIRequest represents the request body for the OpenRouter API.
type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// Message represents a single message in the chat history.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIResponse represents the response body from the OpenRouter API.
type APIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	logger := log.New(os.Stderr, "", 0)

	// Handle version and help flags first
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show this help message")
	flag.Bool("h", false, "Show this help message (alias for --help)")

	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("[ERROR] loading config: %v", err)
	}

	if *showHelp || (flag.Lookup("h") != nil && flag.Lookup("h").Value.String() == "true") {
		printHelp(os.Stdout)
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("book2ru-go v%s\n", version)
		fmt.Printf("Using model: %s\n", cfg.Model)
		os.Exit(0)
	}

	if err := runTranslate(os.Stdin, os.Stdout, logger, cfg); err != nil {
		logger.Fatalf("[ERROR] %v", err)
	}
}

func loadConfig() (*Config, error) {
	// Default configuration
	cfg := &Config{
		Model:          "google/gemini-flash-1.5",
		Prompt:         "Translate this text to Russian. Only return the translated text, nothing else:",
		MetadataFooter: true,
		RetryAttempts:  3,
		StartBatch:     1,
	}

	// Load from .book2ru-go.yml if it exists
	if _, err := os.Stat(".book2ru-go.yml"); err == nil {
		yamlFile, err := os.ReadFile(".book2ru-go.yml")
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
		if err := yaml.Unmarshal(yamlFile, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	// Load from .env file
	_ = godotenv.Load()

	// Parse flags to override config - do NOT use env values as defaults to avoid showing them in help
	flag.StringVar(&cfg.Model, "model", cfg.Model, "Specify LLM model")
	flag.StringVar(&cfg.Model, "m", cfg.Model, "Specify LLM model (alias)")
	flag.StringVar(&cfg.APIKey, "openrouter_key", "", "OpenRouter API key")
	flag.StringVar(&cfg.APIKey, "o", "", "OpenRouter API key (alias)")
	flag.IntVar(&cfg.StartBatch, "start-batch", cfg.StartBatch, "Start translation from specific batch number")

	flag.Parse()

	// Set API key from environment if not set by flag
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("OPENROUTER_KEY")
	}

	return cfg, nil
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage: book2ru-go [options] < input.txt > output.txt")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  book2ru-go < input.txt > output-ru.txt")
	fmt.Fprintln(w, "  book2ru-go --model claude-3-haiku < input.txt > output.txt")
	fmt.Fprintln(w, "  book2ru-go --start-batch 10 < input.txt > output.txt  # Resume from batch 10")
}

func runTranslate(stdin io.Reader, stdout io.Writer, logger *log.Logger, cfg *Config) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("no API key provided. Set OPENROUTER_KEY environment variable or use -o flag")
	}

	if cfg.MetadataFooter {
		logger.Printf("# book2ru-go v%s - starting translation using %s", version, cfg.Model)
	}

	content, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("reading from stdin: %w", err)
	}

	batches := createBatchesFromContent(string(content))
	if cfg.MetadataFooter {
		logger.Printf("# Created %d batches from %d bytes", len(batches), len(content))
		if cfg.StartBatch > 1 {
			logger.Printf("# Starting from batch %d", cfg.StartBatch)
		}
	}

	// Validate start batch
	if cfg.StartBatch < 1 || cfg.StartBatch > len(batches) {
		return fmt.Errorf("invalid start batch %d, must be between 1 and %d", cfg.StartBatch, len(batches))
	}

	// Process batches starting from specified batch
	for i := cfg.StartBatch - 1; i < len(batches); i++ {
		batch := batches[i]
		if cfg.MetadataFooter {
			logger.Printf("# Processing batch %d/%d (%d bytes, %d lines)", i+1, len(batches), batch.Size, batch.Lines)
		}
		
		translated, err := translateTextBatch(batch.Content, cfg, logger)
		if err != nil {
			return fmt.Errorf("translating batch %d: %w\n\nTo resume from this batch, use: --start-batch %d", i+1, err, i+1)
		}
		
		if _, err := fmt.Fprint(stdout, translated); err != nil {
			return fmt.Errorf("writing to stdout: %w", err)
		}
	}

	if cfg.MetadataFooter {
		logger.Printf("# Translated by book2ru-go v%s using model %s", version, cfg.Model)
	}

	return nil
}

// createBatchesFromContent creates batches from content, splitting by 10KB
// This exactly mirrors the Ruby logic
func createBatchesFromContent(content string) []Batch {
	var batches []Batch
	
	// Handle empty content
	if content == "" {
		return batches
	}
	
	// Split content into lines, preserving empty lines at the end like Ruby's split("\n", -1)
	lines := strings.Split(content, "\n")
	
	var currentBatch []string
	currentSize := 0

	for index, line := range lines {
		// Add newline back except for the last line (like Ruby logic)
		var lineWithNewline string
		if index == len(lines)-1 {
			lineWithNewline = line
		} else {
			lineWithNewline = line + "\n"
		}
		
		lineSize := len(lineWithNewline)

		// If adding this line would exceed the limit AND we have content - finish current batch
		if currentSize+lineSize > batchSize && len(currentBatch) > 0 {
			batches = append(batches, Batch{
				Content: strings.Join(currentBatch, ""),
				Size:    currentSize,
				Lines:   len(currentBatch),
			})
			currentBatch = []string{lineWithNewline}
			currentSize = lineSize
		} else {
			currentBatch = append(currentBatch, lineWithNewline)
			currentSize += lineSize
		}
	}

	// Add the last batch if it has content
	if len(currentBatch) > 0 {
		batches = append(batches, Batch{
			Content: strings.Join(currentBatch, ""),
			Size:    currentSize,
			Lines:   len(currentBatch),
		})
	}

	return batches
}

func translateTextBatch(text string, cfg *Config, logger *log.Logger) (string, error) {
	var lastErr error
	
	for attempt := 0; attempt < cfg.RetryAttempts; attempt++ {
		if attempt > 0 {
			logger.Printf("# Retry attempt %d failed: %v", attempt, lastErr)
			time.Sleep(time.Duration(1<<attempt) * time.Second) // Exponential backoff
		}

		apiReq := APIRequest{
			Model: cfg.Model,
			Messages: []Message{
				{Role: "user", Content: fmt.Sprintf("%s\n\n%s", cfg.Prompt, text)},
			},
			Temperature: 0.1,
			Stream:      false,
		}

		reqBody, err := json.Marshal(apiReq)
		if err != nil {
			lastErr = fmt.Errorf("marshalling request body: %w", err)
			continue
		}

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
		if err != nil {
			lastErr = fmt.Errorf("creating http request: %w", err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("HTTP-Referer", httpReferer)
		req.Header.Set("X-Title", xTitle)
		req.Header.Set("User-Agent", userAgent)

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("performing http request: %w", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
			continue
		}

		var apiResp APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			lastErr = fmt.Errorf("decoding API response: %w", err)
			continue // Теперь JSON ошибки тоже ретраятся
		}

		if len(apiResp.Choices) == 0 {
			lastErr = fmt.Errorf("no choices returned from API")
			continue
		}

		content := apiResp.Choices[0].Message.Content
		if cfg.MetadataFooter {
			// Truncate for logging
			logContent := content
			if len(logContent) > 100 {
				logContent = logContent[:100] + "..."
			}
			logger.Printf("# API Response: %s", logContent)
		}
		return content, nil
	}

	return "", fmt.Errorf("failed after %d attempts: %w", cfg.RetryAttempts, lastErr)
} 
