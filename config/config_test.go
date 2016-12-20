package config

import "testing"

func TestLoadConfig(t *testing.T) {
	cf, err := Load("config.toml")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range cf.Servers {
		if k == "db1" && v.Host != "172.16.1.1" {
			t.Fatal("Expected server name db1 not found")
		}
		t.Logf("Key: %s Values: %s\t%s\t%s\t%s\n", k, v.Host, v.Port, v.User, v.Pass)
	}
	mc := len(cf.Monitor)
	if mc != 2 {
		t.Fatal("Expected monitor slice with 2 elements not found, count was ", mc)
	}
	if cf.Monitor[0].Name != "connections" {
		t.Fatal("Did not find expected monitor name connections in element 0")
	}
}
