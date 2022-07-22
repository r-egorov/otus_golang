package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18"` // `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,staff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Tweet struct {
		message string `validate:"len:50"`
		sender  string `validate:"in:user,admin"`
	}

	IntCase struct {
		minimum int `validate:"min:50"`
		maximum int `validate:"max:100"`
		in      int `validate:"in:42"`
	}

	testCase struct {
		name        string
		in          interface{}
		expectedErr error
	}
)

func TestValidateComplex(t *testing.T) {
	tests := []testCase{
		{
			name: "Complex OK",
			in: User{
				ID:     "824e3916-bc8f-4f95-9dd0-d8f62e3b36e7",
				Name:   "John Doe",
				Age:    18,
				Email:  "example@mail.com",
				Role:   "staff",
				Phones: []string{"79234324123", "1111111111"},
				meta:   json.RawMessage(`{"test": "user"}`),
			},
			expectedErr: nil,
		},
		{
			name: "Complex KO",
			in: User{
				ID:     "824e3916-bc8f-4f95-9dd0", // too short
				Name:   "John Doe",
				Age:    10,                // not in set
				Email:  "examplemail.com", // not valid email
				Role:   "meremortal",      // not in
				Phones: []string{"79234324123", "1111111111"},
				meta:   json.RawMessage(`{"test": "user"}`),
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrStringLenInvalid,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrIntUnderMinimum,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrStringNotRegexplike,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrStringNotInSet,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestValidateString(t *testing.T) {
	tests := []testCase{
		{
			name: "String KO - too short",
			in: Tweet{
				message: "Can't you think of even 50 symbols? Pathetic.",
				sender:  "user",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "message",
					Err:   ErrStringLenInvalid,
				},
			},
		},
		{
			name: "String KO - too long",
			in: Tweet{
				message: "Wow, that's a good job you've done here, bravo. You outdid yourself.",
				sender:  "user",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "message",
					Err:   ErrStringLenInvalid,
				},
			},
		},
		{
			name: "String KO - not in set",
			in: Tweet{
				message: "Wow, that's a good job you've done here, bravo. Yo",
				sender:  "justsomedude",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "sender",
					Err:   ErrStringNotInSet,
				},
			},
		},
	}
	runTests(t, tests)
}

func TestValidateInt(t *testing.T) {
	tests := []testCase{
		{
			name: "under minimum",
			in: IntCase{
				minimum: 1,
				maximum: 70,
				in:      42,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "minimum",
					Err:   ErrIntUnderMinimum,
				},
			},
		},
		{
			name: "over maximum",
			in: IntCase{
				minimum: 100,
				maximum: 250,
				in:      42,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "maximum",
					Err:   ErrIntOverMaximum,
				},
			},
		},
		{
			name: "not in set",
			in: IntCase{
				minimum: 100,
				maximum: 100,
				in:      100500,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "in",
					Err:   ErrIntNotInSet,
				},
			},
		},
	}
	runTests(t, tests)
}

func runTests(t *testing.T, testCases []testCase) {
	t.Helper()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			errs := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, errs)
			} else {
				var valErrs ValidationErrors
				if errors.As(errs, &valErrs) {
					var expectedErrs ValidationErrors
					require.ErrorAs(t, tt.expectedErr, &expectedErrs)
					for i, err := range valErrs {
						require.ErrorIs(t, err, expectedErrs[i])
					}
				} else {
					require.ErrorIs(t, errs, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
