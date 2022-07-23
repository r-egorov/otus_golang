package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
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

	StrCase struct {
		Message string   `validate:"len:50"`
		Sender  string   `validate:"in:user,admin"`
		Slice   []string `validate:"regexp:^s\\w*i$"`
	}

	IntCase struct {
		Minimum  int   `validate:"min:50"`
		Maximum  int   `validate:"max:100"`
		In       int   `validate:"in:42"`
		Slice    []int `validate:"in:42,21"`
		Multiple int   `validate:"min:50|max:100"`
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
				Phones: []string{"79234324123", "11111111111"},
				meta:   json.RawMessage(`{"test": "user"}`),
			},
			expectedErr: nil,
		},
		{
			name: "Complex KO",
			in: User{
				ID:     "824e3916-bc8f-4f95-9dd0", // too short
				Name:   "John Doe",
				Age:    10,                                     // not in set
				Email:  "examplemail.com",                      // not valid email
				Role:   "meremortal",                           // not in
				Phones: []string{"792343244123", "1111111111"}, // len not
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
				ValidationError{
					Field: "Phones",
					Err:   ErrStringLenInvalid,
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
			in: StrCase{
				Message: "Can't you think of even 50 symbols? Pathetic.",
				Sender:  "user",
				Slice:   []string{"somesushi", "satoshi", "si"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Message",
					Err:   ErrStringLenInvalid,
				},
			},
		},
		{
			name: "String KO - too long",
			in: StrCase{
				Message: "Wow, that's a good job you've done here, bravo. You outdid yourself.",
				Sender:  "user",
				Slice:   []string{"somesushi", "satoshi", "si"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Message",
					Err:   ErrStringLenInvalid,
				},
			},
		},
		{
			name: "String KO - not in set",
			in: StrCase{
				Message: "Wow, that's a good job you've done here, bravo. Yo",
				Sender:  "justsomedude",
				Slice:   []string{"somesushi", "satoshi", "si"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Sender",
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
			name: "under Minimum",
			in: IntCase{
				Minimum:  1,
				Maximum:  70,
				In:       42,
				Slice:    []int{42, 21},
				Multiple: 70,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Minimum",
					Err:   ErrIntUnderMinimum,
				},
			},
		},
		{
			name: "over Maximum",
			in: IntCase{
				Minimum:  100,
				Maximum:  250,
				In:       42,
				Slice:    []int{42, 21},
				Multiple: 70,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Maximum",
					Err:   ErrIntOverMaximum,
				},
			},
		},
		{
			name: "not in set",
			in: IntCase{
				Minimum:  100,
				Maximum:  100,
				In:       100500,
				Slice:    []int{42, 21},
				Multiple: 70,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "In",
					Err:   ErrIntNotInSet,
				},
			},
		},
		{
			name: "slice not in set",
			in: IntCase{
				Minimum:  100,
				Maximum:  100,
				In:       42,
				Slice:    []int{42, 22},
				Multiple: 70,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Slice",
					Err:   ErrIntNotInSet,
				},
			},
		},
		{
			name: "multiple tag: int under minimum",
			in: IntCase{
				Minimum:  100,
				Maximum:  100,
				In:       42,
				Slice:    []int{42, 21},
				Multiple: 25,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Multiple",
					Err:   ErrIntUnderMinimum,
				},
			},
		},
		{
			name: "multiple tag: int over maximum",
			in: IntCase{
				Minimum:  100,
				Maximum:  100,
				In:       42,
				Slice:    []int{42, 21},
				Multiple: 125,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Multiple",
					Err:   ErrIntOverMaximum,
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
					require.Equal(t, len(expectedErrs), len(valErrs),
						fmt.Sprintf(
							"expected err: %s\ngot err: %s\n",
							expectedErrs, valErrs,
						),
					)
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
