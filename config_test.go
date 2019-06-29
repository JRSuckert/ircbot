package main

import "testing"

func TestParse(t *testing.T) {
	var config Config
	configWant := Config{Name: "Testserver", Adress: "localhost:6667", Nick: "test"}

	testfile := "testconfig.yml"
	config.Parse(testfile)

	if config.Name != configWant.Name {
		t.Errorf("Parse(%s).Name == %s, want %s", testfile, config.Name, configWant.Name)
	}
	if config.Adress != configWant.Adress {
		t.Errorf("Parse(%s).Adress == %s, want %s", testfile, config.Adress, configWant.Adress)
	}
	if config.Nick != configWant.Nick {
		t.Errorf("Parse(%s).Name == %s, want %s", testfile, config.Nick, configWant.Nick)
	}
}
