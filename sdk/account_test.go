package sdk

import "testing"

func TestGenAccount(t *testing.T) {
	sk, pk, addr, err := GenAccount()
	if err != nil {
		t.Fatalf("GenAccount error: %s", err)
	}
	t.Logf("sk: %x\npk: %x\naddr: %x\n", sk, pk, addr)
}
