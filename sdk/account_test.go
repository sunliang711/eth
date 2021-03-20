package sdk

import "testing"

func TestGenAccount(t *testing.T) {
	sk, pk, addr, err := GenAccount()
	if err != nil {
		t.Fatalf("GenAccount error: %s", err)
	}
	t.Logf("sk: %x\npk: %x\naddr: %x\n", sk, pk, addr)
}

func TestExport(t *testing.T) {
	utcFile := "../node/datadirs/node1-datadir/keystore/UTC--2021-03-20T01-48-44.797701000Z--e716f387349d635a3245787174193f104a1759d9"
	bs,err := ExportAccount(utcFile,"sl262732")
	if err != nil{
		t.Fatalf("export account error: %v",err)
	}
	t.Logf("account: %v",string(bs))
}
