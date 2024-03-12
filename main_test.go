package main

import (
	_ "image/png"
	"testing"
)

func Test_createVAO(t *testing.T) {
	type args struct {
		vertices []float32
		indices  []uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createVAO(tt.args.vertices, tt.args.indices); got != tt.want {
				t.Errorf("createVAO() = %v, want %v", got, tt.want)
			}
		})
	}
}
