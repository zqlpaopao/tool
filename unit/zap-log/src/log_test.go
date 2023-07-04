package src
//
//import (
//	"io"
//	"reflect"
//	"testing"
//)
//
//func Test_getWriter(t *testing.T) {
//	type args struct {
//		filename string
//	}
//	tests := []struct {
//		name string
//		args args
//		want io.Writer
//	}{
//		{name: testing.CoverMode(), args: args{filename: testing.CoverMode()}, want: nil},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := getWriter(,); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getWriter() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
