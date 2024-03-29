package auth

import (
	"webdawgengine/database"
	"webdawgengine/identity"
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

func Authenticate(email string, password string) (*identity.Identity, error) {
	// time attack partial mitigation
	// adds up to 0.5 seconds to the response time

	// this technically does not prevent a time attack, since there is still time variance without the randomness added.
	// you could theoretically take an average of a 'valid user; incorrect password' vs 'invalid user' response times
	// to figure out if a user exists, but you would need a lot of data to do that.
	// this should make it *extremely* unlikely to do when paired with 'n login attempt per ip/minute/fingerprint'
	// since you would need way more than `n` login attempts to collect an average

	// https://security.stackexchange.com/questions/96489/can-i-prevent-timing-attacks-with-random-delays/96493#96493
	// https://www.reddit.com/r/PHP/comments/kn6ezp/have_you_secured_your_signup_process_against_a/

	randomSeconds, _ := rand.Int(rand.Reader, big.NewInt(500))
	randomDuration := time.Duration(randomSeconds.Int64()) * time.Millisecond

	time.Sleep(randomDuration)

	id := &identity.Identity{}
	var err error

	user, userErr := database.FetchUserByEmail(email)

	if userErr != nil {
		err = errors.New("invalid email")
		// set user password to dummy password to keep timing consistent when validating password
		user.Password = "$2a$14$KW5OO1wZqGGq3SrpBFj0Oema5DG8Ph7lZJvq0ECkkYBpNFom6b9vO"
	}

	if CheckPassword(password, user.Password) {
		id.IsAuthenticated = true
		id.Email = email

		id.Claims, _ = database.FetchUserClaimsByIdAsMap(user.Id)
	} else {
		err = errors.New("invalid password")
	}

	return id, err
}
