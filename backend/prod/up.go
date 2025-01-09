package prod

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UpProduct struct {
	Ctx context.Context
}

func (u *UpProduct) ConfJsonValidate() (ok bool, err error) {
	result, err := runtime.OpenFileDialog(u.Ctx, runtime.OpenDialogOptions{
		Title: "Open File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Scan json file for tars-release-list",
				Pattern:     "*.json",
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	s, err := os.OpenFile(result, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	b, err := io.ReadAll(s)
	if err != nil {
		fmt.Println(err)
	}
	_, err = json.Marshal(b)
	if err != nil {
		return false, err
	}
	return true, nil
}
