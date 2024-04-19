package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func prepareHttpRequest(token string, method string, endpoint string, bodyData []byte) *http.Request {
	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(bodyData))
	request.Header.Add("Content-Type", "application/json")
	if token != "" {
		request.Header.Add("Authorization", token)
	}
	response := httptest.NewRecorder()
	return
}

func TestViewWorkspaceFiles(t *testing.T) {
	t.Run("returns workspace files", func(t *testing.T) {
		requestString := fmt.Sprintf(`
			{
				"workspace_id": "%s"
			}
			`, "81")
		bodyData := []byte(requestString)

		request, _ := http.NewRequest("POST", "/workspace-files", bytes.NewBuffer(bodyData))
		request.Header.Add("Content-Type", "application/json")

		response := httptest.NewRecorder()

		ViewWorkspaceFiles(response, request)

		got := response.Body
		body, _ := io.ReadAll(got)
		d := make(map[string][]string)
		if err := json.Unmarshal(body, &d); err != nil {
			fmt.Println("Unmarshall Error\n", err)
		}
		want := make(map[string][]string)
		want["file_names"] = []string{"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Mouse_db.sql", "/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Mouse_tinkercad.txt", "/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Mouse_Nonverbal Interaction.pdf"}
		// want["file_names"] = {"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Mouse_db.sql"}

		if !reflect.DeepEqual(want, d) {
			t.Errorf("Got %q, Want %q", d, want)
		}
	})

	t.Run("returns many workspace files", func(t *testing.T) {
		requestString := fmt.Sprintf(`
			{
				"workspace_id": "%s"
			}
			`, "78")
		bodyData := []byte(requestString)

		request, _ := http.NewRequest("POST", "/workspace-files", bytes.NewBuffer(bodyData))
		request.Header.Add("Content-Type", "application/json")

		response := httptest.NewRecorder()

		ViewWorkspaceFiles(response, request)

		got := response.Body
		body, _ := io.ReadAll(got)
		d := make(map[string][]string)
		if err := json.Unmarshal(body, &d); err != nil {
			fmt.Println("Unmarshall Error\n", err)
		}
		want := make(map[string][]string)
		want["file_names"] = []string{
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_filesync.sql",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_write_operation.drawio (5).png",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_[AMV HAJIME NO IPPO] Ippo x Sawamura - _One of Us is Going Down_ (1440x1080) (HDHQ).mp3",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_CALIBRATION 1.0.pdf",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_RM0383 (1).PDF",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_ten_mb.PDF",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_one_mb.docx",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_one.docx",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_time.docx",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_Postman-1699151687493.tar.gz",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_test_a12.tar.gz",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_asd.tar.gz",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_jock_file.exe",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_dev_test.deb",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_lolz.deb",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_kek.deb",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_RM0383 (1).PDF",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_taf.deb",
			"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Dimm_tef2.deb",
		}
		// want["file_names"] = {"/home/aashab/code/src/github.com/aashabtajwar/server-th/storage/Mouse_db.sql"}

		if !reflect.DeepEqual(want, d) {
			t.Errorf("Got %q, Want %q", d, want)
		}
	})

	t.Run("returns 0 workspace files", func(t *testing.T) {
		requestString := fmt.Sprintf(`
			{
				"workspace_id": "%s"
			}
			`, "79")
		bodyData := []byte(requestString)

		request, _ := http.NewRequest("POST", "/workspace-files", bytes.NewBuffer(bodyData))
		request.Header.Add("Content-Type", "application/json")

		response := httptest.NewRecorder()

		ViewWorkspaceFiles(response, request)

		got := response.Body
		body, _ := io.ReadAll(got)
		d := make(map[string][]string)
		if err := json.Unmarshal(body, &d); err != nil {
			fmt.Println("Unmarshall Error\n", err)
		}
		want := 0
		gotten := len(d["file_names"])
		if !reflect.DeepEqual(want, gotten) {
			t.Errorf("Got %q, Want %q", gotten, want)
		}
	})

}
