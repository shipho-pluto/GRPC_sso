package tests

import (
	"grpc_sso/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/shipho-pluto/grpc_proto/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID  = 0
	unrealAppID = int32(2)
	appID       = int32(1)
	appSecret   = "secret"

	passDefaultLen = 10
)

func TestAuth_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)
	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	loginTime := time.Now()

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int32(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestAuth(t *testing.T) {
	test_cases := []struct {
		name  string
		email string
		pass  string
		appID int32
		err   string
	}{
		{
			name:  "normal case",
			email: gofakeit.Email(),
			pass:  randomFakePassword(),
			appID: appID,
		},
		{
			name:  "empty amail",
			email: "",
			pass:  randomFakePassword(),
			appID: appID,
			err:   "email is required",
		},
		{
			name:  "empty amail",
			email: "",
			pass:  "",
			appID: appID,
			err:   "email is required",
		},
		{
			name:  "empty password",
			email: gofakeit.Email(),
			pass:  "",
			appID: appID,
			err:   "password is required",
		},
		{
			name:  "unreal app id",
			email: gofakeit.Email(),
			pass:  randomFakePassword(),
			appID: emptyAppID,
			err:   "app_id is required",
		},
		{
			name:  "unexpected app id",
			email: gofakeit.Email(),
			pass:  randomFakePassword(),
			appID: unrealAppID,
			err:   "unexpected service",
		},
	}

	ctx, st := suite.New(t)

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tc.email,
				Password: tc.pass,
			})
			if err != nil {
				require.Contains(t, err.Error(), tc.err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, respReg.GetUserId())

			respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tc.email,
				Password: tc.pass,
				AppId:    tc.appID,
			})
			if err != nil {
				require.Contains(t, err.Error(), tc.err)
				return
			}

			require.NoError(t, err)
			token := respLog.GetToken()
			require.NotEmpty(t, token)

			tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(appSecret), nil
			})
			require.NoError(t, err)

			loginTime := time.Now()

			claims, ok := tokenParsed.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
			assert.Equal(t, tc.email, claims["email"].(string))
			assert.Equal(t, tc.appID, int32(claims["app_id"].(float64)))

			const deltaSeconds = 1
			assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
		})
	}
}

func TestDoubleRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestLoginWithoutRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.Error(t, err)
	assert.Empty(t, respLog.GetToken())
	assert.ErrorContains(t, err, "invalid credentials")
}

func TestWrongCredentials(t *testing.T) {
	ctx, st := suite.New(t)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: randomFakePassword(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    gofakeit.Email(),
		Password: randomFakePassword(),
		AppId:    appID,
	})
	require.Error(t, err)
	assert.Empty(t, respLog.GetToken())
	assert.ErrorContains(t, err, "invalid credentials")
}
