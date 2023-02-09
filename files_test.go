package gogpt_test

import (
	. "github.com/sashabaranov/go-gpt3"

	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-gpt3/internal/api"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestFileUpload(t *testing.T) {
	api.RegisterHandler("/v1/files", handleCreateFile)
	// create the test server
	var err error
	ts := api.OpenAITestServer()
	ts.Start()
	defer ts.Close()

	client := NewClient(api.TestAPIToken)
	ctx := context.Background()
	client.BaseURL = ts.URL + "/v1"

	req := FileRequest{
		FileName: "test.go",
		FilePath: "api.go",
		Purpose:  "fine-tune",
	}
	_, err = client.CreateFile(ctx, req)
	if err != nil {
		t.Fatalf("CreateFile error: %v", err)
	}
}

// handleCreateFile Handles the images endpoint by the test server.
func handleCreateFile(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// edits only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err = r.ParseMultipartForm(1024 * 1024 * 1024)
	if err != nil {
		http.Error(w, "file is more than 1GB", http.StatusInternalServerError)
		return
	}

	values := r.Form
	var purpose string
	for key, value := range values {
		if key == "purpose" {
			purpose = value[0]
		}
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return
	}
	defer file.Close()

	var fileReq = File{
		Bytes:     int(header.Size),
		ID:        strconv.Itoa(int(time.Now().Unix())),
		FileName:  header.Filename,
		Purpose:   purpose,
		CreatedAt: time.Now().Unix(),
		Object:    "test-objecct",
		Owner:     "test-owner",
	}

	resBytes, _ = json.Marshal(fileReq)
	fmt.Fprint(w, string(resBytes))

}