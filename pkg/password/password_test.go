package password

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "Test123!@#",
			wantErr:  nil,
		},
		{
			name:     "password too short",
			password: "Test1!",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "password too long",
			password: "Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!Test1!",
			wantErr:  ErrPasswordTooLong,
		},
		{
			name:     "missing uppercase",
			password: "test123!@#",
			wantErr:  ErrMissingUpper,
		},
		{
			name:     "missing lowercase",
			password: "TEST123!@#",
			wantErr:  ErrMissingLower,
		},
		{
			name:     "missing number",
			password: "TestTest!@#",
			wantErr:  ErrMissingNumber,
		},
		{
			name:     "missing special",
			password: "Test12345",
			wantErr:  ErrMissingSpecial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.password)
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashAndVerify(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "Test123!@#",
		},
		{
			name:     "complex password",
			password: "VeryC0mplex!@#$%^&*()",
		},
		{
			name:     "long password",
			password: "ThisIsAVeryLongPassword123!@#WithLotsOfCharacters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test hashing
			hashedPassword, err := Hash(tt.password)
			if err != nil {
				t.Errorf("Hash() error = %v", err)
				return
			}
			if hashedPassword == tt.password {
				t.Error("Hash() returned original password")
			}

			// Test verification
			if !Verify(tt.password, hashedPassword) {
				t.Error("Verify() failed to verify valid password")
			}

			// Test verification with wrong password
			if Verify("wrong"+tt.password, hashedPassword) {
				t.Error("Verify() verified invalid password")
			}
		})
	}
}
