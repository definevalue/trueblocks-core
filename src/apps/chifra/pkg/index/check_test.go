package chunk

import (
	"reflect"
	"testing"
)

func TestGetAddress(t *testing.T) {
	type args struct {
		chunk  Chunk
		number int
	}
	tests := []struct {
		name    string
		args    args
		want    AddressRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddress(tt.args.chunk, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
