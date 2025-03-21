package tests

import (
	"kerberos/as"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {

	checkUser := func(t testing.TB, username, password string, want bool) {
		got := as.AuthenticateUser(username, password)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}

	checkUser(t, "kerberos", "Password", true)
	checkUser(t, "kerberos1", "Password1", true)

}
