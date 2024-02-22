//go:generate stringer -type Type

package milestone

import "testing"

func Test_typeFromString(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "empty string",
			args: args{
				t: "",
			},
			want: unknown,
		},
		{
			name: "not defined",
			args: args{
				t: "asdfasfdasdf",
			},
			want: unknown,
		},
		{
			name: "first",
			args: args{
				t: Type(1).String(),
			},
			want: Type(1),
		},
		{
			name: "last",
			args: args{
				t: Type(typesCount - 1).String(),
			},
			want: typesCount - 1,
		},
		{
			name: "instance deleted",
			args: args{
				t: InstanceDeleted.String(),
			},
			want: InstanceDeleted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeFromString(tt.args.t); got != tt.want {
				t.Errorf("typeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
