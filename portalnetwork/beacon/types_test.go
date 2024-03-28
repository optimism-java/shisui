package beacon

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestForkedLightClientUpdateRange(t *testing.T) {
	filePath := "testdata/light_client_updates_by_range.json"

	f, _ := os.Open(filePath)
	jsonStr, _ := io.ReadAll(f)

	var result map[string]interface{}
	_ = json.Unmarshal(jsonStr, &result)

	fmt.Println(result)
}
