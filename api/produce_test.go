package api

import "testing"

func Test_isValidProduceCode(t *testing.T) {
	type args struct {
		produceCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid product code",
			args: args{
				produceCode: "A12T-4GH7-QPL9-3N4M",
			},
			want: true,
		},
		{
			name: "valid product code is case insensitive",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-3n4m",
			},
			want: true,
		},
		{
			name: "invalid product code only three sets",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-",
			},
			want: false,
		},
		{
			name: "invalid product code special character !",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-ALL!",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidProduceCode(tt.args.produceCode); got != tt.want {
				t.Errorf("isValidProduceCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
