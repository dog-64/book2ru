package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBatchesFromContent(t *testing.T) {
	// Save original batchSize and restore it after tests
	originalBatchSize := batchSize
	defer func() { batchSize = originalBatchSize }()

	tests := []struct {
		name        string
		content     string
		batchSize   int
		wantBatches int
		wantContent []string
	}{
		{
			name:        "Empty content",
			content:     "",
			batchSize:   100,
			wantBatches: 0,
			wantContent: []string{},
		},
		{
			name:        "Single small batch",
			content:     "hello world",
			batchSize:   100,
			wantBatches: 1,
			wantContent: []string{"hello world"},
		},
		{
			name:        "Multiple batches",
			content:     "This is the first line.\nThis is the second line, which is much longer and will cause a new batch because it exceeds the limit.",
			batchSize:   100,
			wantBatches: 2,
			wantContent: []string{
				"This is the first line.\n",
				"This is the second line, which is much longer and will cause a new batch because it exceeds the limit.",
			},
		},
		{
			name:        "Content with trailing newline",
			content:     "line1\nline2\n",
			batchSize:   100,
			wantBatches: 1,
			wantContent: []string{"line1\nline2\n"},
		},
		{
			name:        "Content exactly at batch size",
			content:     strings.Repeat("a", 50) + "\n" + strings.Repeat("b", 49),
			batchSize:   100,
			wantBatches: 1,
			wantContent: []string{strings.Repeat("a", 50) + "\n" + strings.Repeat("b", 49)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchSize = tt.batchSize
			batches := createBatchesFromContent(tt.content)
			
			if len(batches) != tt.wantBatches {
				t.Errorf("createBatchesFromContent() got %d batches, want %d", len(batches), tt.wantBatches)
			}
			
			for i, batch := range batches {
				if i >= len(tt.wantContent) {
					t.Fatalf("got more batches than expected content")
				}
				if batch.Content != tt.wantContent[i] {
					t.Errorf("batch %d content mismatch:\ngot:  %q\nwant: %q", i, batch.Content, tt.wantContent[i])
				}
			}
		})
	}
}

func TestRunTranslate(t *testing.T) {
	// Save original apiURL and restore it after tests
	originalAPIURL := apiURL
	defer func() { apiURL = originalAPIURL }()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := APIResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{Content: "Привет, мир"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	apiURL = server.URL

	t.Run("successful translation", func(t *testing.T) {
		cfg := &Config{
			Model:          "test-model",
			Prompt:         "test-prompt",
			MetadataFooter: true,
			RetryAttempts:  1,
			StartBatch:     1,
			APIKey:         "test-key",
		}

		stdin := strings.NewReader("hello world")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		logger := log.New(&stderr, "", 0)

		err := runTranslate(stdin, &stdout, logger, cfg)
		if err != nil {
			t.Fatalf("runTranslate() returned error: %v", err)
		}

		if got := stdout.String(); got != "Привет, мир" {
			t.Errorf("stdout = %q; want %q", got, "Привет, мир")
		}
		if !strings.Contains(stderr.String(), "starting translation") {
			t.Errorf("stderr should contain starting message, got: %q", stderr.String())
		}
	})

	t.Run("start from specific batch", func(t *testing.T) {
		// Save original batchSize and restore it after test
		originalBatchSize := batchSize
		defer func() { batchSize = originalBatchSize }()
		batchSize = 3 // Small batch size to create multiple batches

		cfg := &Config{
			Model:          "test-model",
			Prompt:         "test-prompt",
			MetadataFooter: true,
			RetryAttempts:  1,
			StartBatch:     2, // Start from second batch
			APIKey:         "test-key",
		}

		// Content that will create 3 batches with batchSize=3: "aa\n", "bb\n", "cc"
		stdin := strings.NewReader("aa\nbb\ncc")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		logger := log.New(&stderr, "", 0)

		err := runTranslate(stdin, &stdout, logger, cfg)
		if err != nil {
			t.Fatalf("runTranslate() returned error: %v", err)
		}

		// Should process only batches 2 and 3, so 2 API calls = 2 outputs
		expectedOutput := "Привет, мирПривет, мир"
		if got := stdout.String(); got != expectedOutput {
			t.Errorf("stdout = %q; want %q", got, expectedOutput)
		}
		
		stderrStr := stderr.String()
		if !strings.Contains(stderrStr, "Starting from batch 2") {
			t.Errorf("stderr should contain start batch message, got: %q", stderrStr)
		}
		if !strings.Contains(stderrStr, "Processing batch 2/3") {
			t.Errorf("stderr should show processing batch 2, got: %q", stderrStr)
		}
		if !strings.Contains(stderrStr, "Processing batch 3/3") {
			t.Errorf("stderr should show processing batch 3, got: %q", stderrStr)
		}
		if strings.Contains(stderrStr, "Processing batch 1/3") {
			t.Errorf("stderr should NOT show processing batch 1, got: %q", stderrStr)
		}
	})

	t.Run("error message includes resume instruction", func(t *testing.T) {
		// Save original batchSize and restore it after test
		originalBatchSize := batchSize
		defer func() { batchSize = originalBatchSize }()
		batchSize = 3 // Small batch size to create multiple batches

		// Mock server that always returns error
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer errorServer.Close()

		// Temporarily override apiURL for this test
		originalAPIURL := apiURL
		apiURL = errorServer.URL
		defer func() { apiURL = originalAPIURL }()

		cfg := &Config{
			Model:          "test-model",
			Prompt:         "test-prompt",
			MetadataFooter: false, // Disable metadata to focus on error message
			RetryAttempts:  1,     // Only 1 retry to speed up test
			StartBatch:     1,
			APIKey:         "test-key",
		}

		// Content that will create 3 batches, error should happen on batch 1
		stdin := strings.NewReader("aa\nbb\ncc")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		logger := log.New(&stderr, "", 0)

		err := runTranslate(stdin, &stdout, logger, cfg)
		if err == nil {
			t.Fatal("runTranslate() should have returned error")
		}

		errorMsg := err.Error()
		if !strings.Contains(errorMsg, "translating batch 1:") {
			t.Errorf("error should mention batch 1, got: %q", errorMsg)
		}
		if !strings.Contains(errorMsg, "To resume from this batch, use: --start-batch 1") {
			t.Errorf("error should include resume instruction, got: %q", errorMsg)
		}
	})

	t.Run("invalid start batch", func(t *testing.T) {
		cfg := &Config{
			Model:          "test-model",
			StartBatch:     10, // Invalid batch number
			APIKey:         "test-key",
		}

		stdin := strings.NewReader("hello world") // Only creates 1 batch
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		logger := log.New(&stderr, "", 0)

		err := runTranslate(stdin, &stdout, logger, cfg)
		if err == nil || !strings.Contains(err.Error(), "invalid start batch 10, must be between 1 and 1") {
			t.Errorf("Expected invalid start batch error, got: %v", err)
		}
	})

	t.Run("missing api key", func(t *testing.T) {
		cfg := &Config{APIKey: ""}
		err := runTranslate(nil, nil, log.New(io.Discard, "", 0), cfg)
		if err == nil || !strings.Contains(err.Error(), "no API key provided") {
			t.Errorf("Expected 'no API key' error, got: %v", err)
		}
	})
}

func TestTranslateTextBatch_Retry(t *testing.T) {
	// Save original apiURL and restore it after tests
	originalAPIURL := apiURL
	defer func() { apiURL = originalAPIURL }()

	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"success"}}]}`))
	}))
	defer server.Close()

	apiURL = server.URL

	cfg := &Config{
		Model:         "test-model",
		Prompt:        "test",
		RetryAttempts: 3,
		StartBatch:    1,
		APIKey:        "key",
	}

	var stderr bytes.Buffer
	logger := log.New(&stderr, "", 0)

	result, err := translateTextBatch("test text", cfg, logger)

	if err != nil {
		t.Fatalf("translateTextBatch() failed: %v", err)
	}
	if result != "success" {
		t.Errorf("expected 'success', got %q", result)
	}
	if attempt != 3 {
		t.Errorf("expected 3 attempts, got %d", attempt)
	}
	if !strings.Contains(stderr.String(), "Retry attempt 1") || !strings.Contains(stderr.String(), "Retry attempt 2") {
		t.Errorf("stderr should show retry attempts, got: %q", stderr.String())
	}
}

func TestCreateBatchesFromContent_EdgeCases(t *testing.T) {
	originalBatchSize := batchSize
	defer func() { batchSize = originalBatchSize }()
	batchSize = 10

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "Single line longer than batch size",
			content: "this is a very long line that exceeds the batch size",
			want:    []string{"this is a very long line that exceeds the batch size"},
		},
		{
			name:    "Empty lines",
			content: "\n\n\n",
			want:    []string{"\n\n\n"},
		},
		{
			name:    "Mixed content",
			content: "short\nvery long line that definitely exceeds our small batch size\nshort again",
			want: []string{
				"short\n",
				"very long line that definitely exceeds our small batch size\n",
				"short again",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := createBatchesFromContent(tt.content)
			got := make([]string, len(batches))
			for i, batch := range batches {
				got[i] = batch.Content
			}

			if len(got) != len(tt.want) {
				t.Errorf("got %d batches, want %d", len(got), len(tt.want))
			}

			for i, content := range got {
				if i >= len(tt.want) {
					break
				}
				if content != tt.want[i] {
					t.Errorf("batch %d: got %q, want %q", i, content, tt.want[i])
				}
			}
		})
	}
} 