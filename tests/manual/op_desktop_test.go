package onepasswordmanager

import (
	"testing"

	opentofu "github.com/freethebaerchen/tpmsfe/tests"
)

func TestOnepasswordDesktop(t *testing.T) {
	projectPath := "/Users/jochen/Projects/tpmsfe/opentofu-test-project/password_managers/onepassword/desktop"
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
