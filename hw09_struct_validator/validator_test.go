package hw09structvalidator

import (
	"encoding/json"
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
		Role   UserRole `validate:"in:admin,stuff"`
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
)

var (
	ValidUser = User{
		ID:     "1295bf4e-7f04-49b7-8a91-ec94c2803f1d",
		Name:   "Vally",
		Age:    28,
		Email:  "amazing@coder.com",
		Role:   "admin",
		Phones: []string{"+1234567891", "+4232517842"},
		meta:   nil,
	}
	WrongUser = User{
		ID:     "13",
		Name:   "Wrongy",
		Age:    53,
		Email:  "soo_wrong@gmail_com",
		Role:   "staff",
		Phones: []string{"+123456789123", "14123"},
		meta:   nil,
	}

	Uber    = App{Version: "4.371.10007"}
	MindMap = App{Version: "6.122"}

	CryptoToken = Token{
		Header:    nil,
		Payload:   nil,
		Signature: nil,
	}

	WebResponse = Response{
		Code: 200,
		Body: `PUT /files/129742 HTTP/1.1\r\n
Host: example.com\r\n
User-Agent: Chrome/54.0.2803.1\r\n
Content-Length: 202\r\n
\r\n
This is a message body. All content in this message body should be stored under the 
/files/129742 path, as specified by the PUT specification. The message body does
not have to be terminated with CRLF.`,
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: ValidUser, expectedErr: nil},
		{in: WrongUser, expectedErr: ValidationErrors{
			ValidationError{Field: "ID", Err: ErrFailedLen},
			ValidationError{Field: "Age", Err: ErrFailedMinMax},
			ValidationError{Field: "Email", Err: ErrFailedRegexp},
			ValidationError{Field: "Role", Err: ErrFailedIn},
			ValidationError{Field: "Phones", Err: ErrFailedLen},
		}},
		{in: Uber, expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrFailedLen}}},
		{in: MindMap, expectedErr: nil},
		{in: WebResponse, expectedErr: nil},
		{in: CryptoToken, expectedErr: nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errors := Validate(tt.in)

			validationErrs, ok := errors.(ValidationErrors)
			if !ok {
				require.Equal(t, tt.expectedErr, errors)
			} else {
				expectedValidationErrs, ok := tt.expectedErr.(ValidationErrors)
				require.True(t, ok)

				for i := 0; i < len(expectedValidationErrs); i++ {
					require.ErrorIs(t, validationErrs[i].Err, expectedValidationErrs[i].Err)
				}
			}
		})
	}
}
