package config

import (
	"testing"

	"github.com/tlmiller/bookie/pkg/executor"
)

func TestRegisterMakeDoubleKey(t *testing.T) {
	maker := Maker(func(c map[string]string) (executor.Executor, error) {
		return nil, nil
	})
	err := RegisterMaker("test", maker)
	if err != nil {
		t.Fatalf("unexpected error when adding unqiue maker: %v", err)
	}
	err = RegisterMaker("test", maker)
	if err == nil {
		t.Fatal("should have recieved error for double register of maker")
	}
}

func TestMakeExcutorNoMaker(t *testing.T) {
	e := executorConfig{Type: "noexist"}
	_, err := makeExecutor(&e)
	if err == nil {
		t.Fatal("should have recieved error for no maker")
	}
}
