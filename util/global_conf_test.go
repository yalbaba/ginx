package util

import "testing"

func TestGlobalConf_Reload(t *testing.T) {
	g := GlobalConf{}
	g.Reload()
	t.Logf("%+v", g)
}
