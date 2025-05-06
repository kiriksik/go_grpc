package tests

import (
	grpс_auth "grpc_auth"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kiriksik/go_grpc/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "abobus"
	passDefoultLen = 10
)

func TestRegister_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &grpс_auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := st.AuthClient.Login(ctx, &grpс_auth.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, tokenParsed)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &grpс_auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg_2, err := st.AuthClient.Register(ctx, &grpс_auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg_2.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefoultLen)
}
