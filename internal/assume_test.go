package internal

import (
	"flag"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/aws/aws-sdk-go/service/sso/ssoiface"
	"github.com/aws/aws-sdk-go/service/ssooidc"
	"github.com/aws/aws-sdk-go/service/ssooidc/ssooidciface"
	. "github.com/theurichde/go-aws-sso/pkg/sso"
	"github.com/urfave/cli/v2"
)

type mockSSOOIDCClient struct {
	t *testing.T
	ssooidciface.SSOOIDCAPI
	CreateTokenOutput              ssooidc.CreateTokenOutput
	RegisterClientOutput           ssooidc.RegisterClientOutput
	StartDeviceAuthorizationOutput ssooidc.StartDeviceAuthorizationOutput
}

func (m mockSSOOIDCClient) CreateToken(i *ssooidc.CreateTokenInput) (*ssooidc.CreateTokenOutput, error) {
	m.t.Log(i)
	return &m.CreateTokenOutput, nil
}

func (m mockSSOOIDCClient) StartDeviceAuthorization(i *ssooidc.StartDeviceAuthorizationInput) (*ssooidc.StartDeviceAuthorizationOutput, error) {
	m.t.Log(i)
	return &m.StartDeviceAuthorizationOutput, nil
}

func (m mockSSOOIDCClient) RegisterClient(i *ssooidc.RegisterClientInput) (*ssooidc.RegisterClientOutput, error) {
	m.t.Log(i)
	return &m.RegisterClientOutput, nil
}

type mockSSOClient struct {
	t *testing.T
	ssoiface.SSOAPI
	GetRoleCredentialsOutput sso.GetRoleCredentialsOutput
	ListAccountRolesOutput   sso.ListAccountRolesOutput
	ListAccountsOutput       sso.ListAccountsOutput
}

func (m mockSSOClient) GetRoleCredentials(i *sso.GetRoleCredentialsInput) (*sso.GetRoleCredentialsOutput, error) {
	m.t.Log(i)
	return &m.GetRoleCredentialsOutput, nil
}

func TestAssumeDirectly(t *testing.T) {
	t.Log("TestAssumeDirectly")
	temp, err := os.CreateTemp("", "go-aws-sso-assume-directly_")
	if err != nil {
		t.Error(err)
	}
	CredentialsFilePath = temp.Name()
	defer func(path string) {
		t.Log("wat", path)
		if r := recover(); r != nil {
			t.Log("Recovered in f", r)
		}
		err := os.RemoveAll(path)
		if err != nil {
			t.Error(err)
		}
	}(CredentialsFilePath)
	t.Log("TestAssumeDirectly", CredentialsFilePath)

	dummyInt := int64(132465)
	dummy := "dummy_assume_directly"
	accessToken := "AccessToken"

	ssoClient := mockSSOClient{
		t:      t,
		SSOAPI: nil,
		GetRoleCredentialsOutput: sso.GetRoleCredentialsOutput{RoleCredentials: &sso.RoleCredentials{
			AccessKeyId:     &dummy,
			Expiration:      &dummyInt,
			SecretAccessKey: &dummy,
			SessionToken:    &dummy,
		}},
	}
	t.Log("TestAssumeDirectly")

	expires := int64(0)

	oidcClient := mockSSOOIDCClient{
		t:          t,
		SSOOIDCAPI: nil,
		CreateTokenOutput: ssooidc.CreateTokenOutput{
			AccessToken: &accessToken,
		},
		RegisterClientOutput: ssooidc.RegisterClientOutput{
			AuthorizationEndpoint: &dummy,
			ClientId:              &dummy,
			ClientSecret:          &dummy,
			ClientSecretExpiresAt: &expires,
			TokenEndpoint:         &dummy,
		},
		StartDeviceAuthorizationOutput: ssooidc.StartDeviceAuthorizationOutput{
			DeviceCode:              &dummy,
			UserCode:                &dummy,
			VerificationUri:         &dummy,
			VerificationUriComplete: &dummy,
		},
	}

	t.Log("TestAssumeDirectly")
	flagSet := flag.NewFlagSet("test-set", flag.ContinueOnError)
	flagSet.String("start-url", "foobar", "")
	flagSet.String("region", "eu-central-1", "")
	flagSet.String("account-id", "123456", "")
	flagSet.String("role-name", "super-admin", "")
	flagSet.String("profile", "default", "")
	flagSet.Bool("persist", true, "")
	ctx := cli.NewContext(nil, flagSet, nil)

	t.Log("TestAssumeDirectly")
	AssumeDirectly(oidcClient, ssoClient, ctx)
	t.Log("TestAssumeDirectly")

	content, err := os.ReadFile(CredentialsFilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log("TestAssumeDirectly", CredentialsFilePath)
	t.Log("TestAssumeDirectly")
	got := string(content)
	t.Log("TestAssumeDirectly")
	want := "[default]\naws_access_key_id     = dummy_assume_directly\naws_secret_access_key = dummy_assume_directly\naws_session_token     = dummy_assume_directly\nregion                = eu-central-1\n"

	t.Log("TestAssumeDirectly")
	if got != want {
		t.Errorf("Got: %v, but wanted: %v", got, want)
	}
	t.Log("TestAssumeDirectly")

}
