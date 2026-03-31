package mask_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	t.Parallel()

	type args struct {
		pass string
	}

	type want struct {
		masked string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty password",
			args: args{
				pass: "",
			},
			want: want{
				masked: "*",
			},
		},
		{
			name: "password",
			args: args{
				pass: "password",
			},
			want: want{
				masked: "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mask.PasswordMask(tt.args.pass)

			assert.Equalf(t, tt.want.masked, got,
				"Password() = %v, want %v", got, tt.want)
		})
	}
}

//nolint:funlen
func TestEmail(t *testing.T) {
	t.Parallel()

	type args struct {
		email string
	}

	type want struct {
		masked string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty email",
			args: args{
				email: "",
			},
			want: want{
				masked: "",
			},
		},
		{
			name: "correct email",
			args: args{
				email: "correct@email.com",
			},
			want: want{
				masked: "c*****t@email.com",
			},
		},
		{
			name: "mini email",
			args: args{
				email: "c2@email.com",
			},
			want: want{
				masked: "c2@email.com",
			},
		},
		{
			name: "incorrect email",
			args: args{
				email: "incorrect.email.com",
			},
			want: want{
				masked: "incorrect.email.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mask.EmailMask(tt.args.email)

			assert.Equalf(t, tt.want.masked, got,
				"Email() = %v, want %v", got, tt.want)
		})
	}
}

//nolint:funlen
func TestPhone(t *testing.T) {
	t.Parallel()

	type args struct {
		phone string
	}

	type want struct {
		masked string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty phone",
			args: args{
				phone: "",
			},
			want: want{
				masked: "",
			},
		},
		{
			name: "russian phone",
			args: args{
				phone: "+78005553535",
			},
			want: want{
				masked: "+7800*****35",
			},
		},
		{
			name: "not russian phone",
			args: args{
				phone: "+7858005553535",
			},
			want: want{
				masked: "+785800*****35",
			},
		},
		{
			name: "mini phone",
			args: args{
				phone: "8005553535",
			},
			want: want{
				masked: "800*****35",
			},
		},
		{
			name: "incorrect phone",
			args: args{
				phone: "553535",
			},
			want: want{
				masked: "****35",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mask.PhoneMask(tt.args.phone)

			assert.Equalf(t, tt.want.masked, got,
				"Phone() = %v, want %v", got, tt.want)
		})
	}
}
