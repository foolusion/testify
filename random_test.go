package main

import (
	"fmt"
	"net/url"
	"testing"
)

func BenchmarkMapJoin(b *testing.B) {
	vals := url.Values(map[string][]string{
		"c": []string{"1", "2"},
	})
	for i := 0; i < b.N; i++ {
		mapJoin(vals, ":", "=")
	}
}

func TestMapJoin(t *testing.T) {
	tests := []struct {
		in   url.Values
		want string
	}{
		{
			in:   url.Values(map[string][]string{"c": []string{"1", "2"}, "b": []string{"3", "4"}, "a": []string{"5", "6"}}),
			want: "a=5,6:b=3,4:c=1,2",
		},
		{
			in:   url.Values(map[string][]string{"userid": []string{"123", "456"}, "test": []string{"poop"}}),
			want: "test=poop:userid=123,456",
		},
	}
	for _, v := range tests {
		out := mapJoin(v.in, ":", "=")
		if out != v.want {
			t.Errorf("got mapJoin(%v, %q, %q) = %v want %v", v.in, ":", "=", out, v.want)
		}
	}
}

func TestGetUniform(t *testing.T) {
	tests := []struct {
		hash     hashedUnit
		min, max float64
		want     float64
	}{
		{
			hash: 0,
			min:  0,
			max:  1,
			want: 0,
		},
		{
			hash: 0xFFFFFFFFFFFFFFF,
			min:  0,
			max:  16.1,
			want: 16.1,
		},
	}

	for _, v := range tests {
		out := v.hash.getUniform(v.min, v.max)
		if out != v.want {
			t.Errorf("got (%v).getUniform(%v, %v) = %v want %v", v.hash, v.min, v.max, out, v.want)
		}
	}
}

func TestRandomInt(t *testing.T) {
	tests := []struct {
		hash     hashedUnit
		min, max int64
		want     int64
	}{
		{hash: 0, min: 0, max: 1, want: 0},
		{hash: 0xFFFFFFFFFFFFFFE, min: 0, max: 10, want: 10},
	}
	for _, v := range tests {
		out := v.hash.randomInt(v.min, v.max)
		if out != v.want {
			t.Errorf("got (%v).randomInt(%v, %v) = %v want %v", v.hash, v.min, v.max, out, v.want)
		}
	}
}

func TestBernoulliTrial(t *testing.T) {
	tests := []struct {
		hash    hashedUnit
		p       float64
		wantB   bool
		wantErr error
	}{
		{hash: 0, p: .5, wantB: false, wantErr: nil},
		{hash: 0xFFFFFFFFFFFFFFF, p: .5, wantB: true, wantErr: nil},
		{hash: 0, p: 1.2, wantB: false, wantErr: fmt.Errorf("p must be between 0 and 1: %v", 1.2)},
		{hash: 0, p: -1, wantB: false, wantErr: fmt.Errorf("p must be between 0 and 1: %v", -1)},
	}

	for _, v := range tests {
		out, err := v.hash.bernoulliTrial(v.p)
		if out != v.wantB {
			t.Errorf("got (%v).bernoulliTrial(%v) = %v want %v", v.hash, v.p, out, v.wantB)
		}
		if err != nil && v.wantErr != nil && err.Error() != v.wantErr.Error() {
			t.Errorf("got (%v).bernoulliTrial(%v) = %v want %v", v.hash, v.p, err, v.wantErr)
		}
	}
}
