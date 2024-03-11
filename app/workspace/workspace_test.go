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

}
