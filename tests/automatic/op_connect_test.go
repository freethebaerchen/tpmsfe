package automatictests

import (
	"os"
	"testing"

	opentofu "github.com/freethebaerchen/tpmsfe/tests"
)

func TestOnepasswordConnect(t *testing.T) {
	currentDir := os.Getenv("PWD")
	projectPath := currentDir + "/../../opentofu-test-project/password_managers/onepassword/connect"
	outputs, err := opentofu.TofuApply(projectPath)
	if err != nil {
		t.Fatalf("TofuApply failed: %v", err)
	}
	encrypted, err := opentofu.CheckStateForRandValue(projectPath, "terraform.tfstate", outputs)
	if err != nil {
		t.Fatalf("CheckStateForRandValue failed: %v", err)
	}

	if !encrypted {
		t.Fatalf("The state encryption did not work. It is still in plain text.")
	}
}
