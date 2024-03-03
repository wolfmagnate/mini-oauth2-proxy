package headerInjection

import (
	"fmt"

	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/session"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
)

type headerInjector interface {
	GetKey() string
	GetValue(id sessionid.ID) (string, error)
}

type idTokenInjector struct {
	Name   string
	Claims []ClaimType
}

type userInfoInjector struct {
	Name   string
	Claims []ClaimType
}

func (injector *idTokenInjector) GetKey() string {
	return injector.Name
}

func (injector *idTokenInjector) GetValue(id sessionid.ID) (string, error) {
	idToken, err := session.GetIDToken(id)
	if err != nil {
		return "", err
	}
	var tokenClaims claims
	if err := idToken.Claims(&tokenClaims); err != nil {
		return "", err
	}
	for _, claim := range injector.Claims {
		if tokenClaims.has(claim) {
			return tokenClaims.get(claim)
		}
	}
	return "", fmt.Errorf("error: no valid claim found to set %v", injector.Name)
}

func (injector *userInfoInjector) GetKey() string {
	return injector.Name
}

func (injector *userInfoInjector) GetValue(id sessionid.ID) (string, error) {
	userinfo, err := session.GetUserInfo(id)
	if err != nil {
		return "", err
	}
	var userinfoClaims claims
	if err := userinfo.Claims(&userinfoClaims); err != nil {
		return "", err
	}
	for _, claim := range injector.Claims {
		if userinfoClaims.has(claim) {
			return userinfoClaims.get(claim)
		}
	}
	return "", fmt.Errorf("error: no valid claim found to set %v", injector.Name)
}

type claims struct {
	Name                string `json:"name"`
	FamilyName          string `json:"family_name"`
	GivenName           string `json:"given_name"`
	MiddleName          string `json:"middle_name"`
	Nickname            string `json:"nickname"`
	PreferredUsername   string `json:"preferred_username"`
	Profile             string `json:"profile"`
	Picture             string `json:"picture"`
	Website             string `json:"website"`
	Gender              string `json:"gender"`
	Birthdate           string `json:"birthdate"`
	Zoneinfo            string `json:"zoneinfo"`
	Locale              string `json:"locale"`
	UpdatedAt           string `json:"updated_at"`
	Email               string `json:"email"`
	EmailVerified       *bool  `json:"email_verified"`
	Address             string `json:"address"`
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified *bool  `json:"phone_number_verified"`
}

func (c *claims) has(claimType ClaimType) bool {
	switch claimType {
	case Name:
		return c.Name != ""
	case FamilyName:
		return c.FamilyName != ""
	case GivenName:
		return c.GivenName != ""
	case MiddleName:
		return c.MiddleName != ""
	case Nickname:
		return c.Nickname != ""
	case PreferredUsername:
		return c.PreferredUsername != ""
	case Profile:
		return c.Profile != ""
	case Picture:
		return c.Picture != ""
	case Website:
		return c.Website != ""
	case Gender:
		return c.Gender != ""
	case Birthdate:
		return c.Birthdate != ""
	case Zoneinfo:
		return c.Zoneinfo != ""
	case Locale:
		return c.Locale != ""
	case UpdatedAt:
		return c.UpdatedAt != ""
	case Email:
		return c.Email != ""
	case EmailVerified:
		return c.EmailVerified != nil
	case Address:
		return c.Address != ""
	case PhoneNumber:
		return c.PhoneNumber != ""
	case PhoneNumberVerified:
		return c.PhoneNumberVerified != nil
	default:
		return false
	}
}

func (c *claims) get(claimType ClaimType) (string, error) {
	switch claimType {
	case Name:
		return c.Name, nil
	case FamilyName:
		return c.FamilyName, nil
	case GivenName:
		return c.GivenName, nil
	case MiddleName:
		return c.MiddleName, nil
	case Nickname:
		return c.Nickname, nil
	case PreferredUsername:
		return c.PreferredUsername, nil
	case Profile:
		return c.Profile, nil
	case Picture:
		return c.Picture, nil
	case Website:
		return c.Website, nil
	case Gender:
		return c.Gender, nil
	case Birthdate:
		return c.Birthdate, nil
	case Zoneinfo:
		return c.Zoneinfo, nil
	case Locale:
		return c.Locale, nil
	case UpdatedAt:
		return c.UpdatedAt, nil
	case Email:
		return c.Email, nil
	case EmailVerified:
		return getStringOfBool(*c.EmailVerified), nil
	case Address:
		return c.Address, nil
	case PhoneNumber:
		return c.PhoneNumber, nil
	case PhoneNumberVerified:
		return getStringOfBool(*c.PhoneNumberVerified), nil
	default:
		return "", fmt.Errorf("error: claim type %v not found", claimType)
	}
}

func getStringOfBool(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
