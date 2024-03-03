package headerInjection

import "errors"

type ClaimType int

const (
	Name ClaimType = iota
	FamilyName
	GivenName
	MiddleName
	Nickname
	PreferredUsername
	Profile
	Picture
	Website
	Gender
	Birthdate
	Zoneinfo
	Locale
	UpdatedAt
	Email
	EmailVerified
	Address
	PhoneNumber
	PhoneNumberVerified
)

func toClaimType(s string) (ClaimType, error) {
	switch s {
	case "name":
		return Name, nil
	case "family_name":
		return FamilyName, nil
	case "given_name":
		return GivenName, nil
	case "middle_name":
		return MiddleName, nil
	case "nickname":
		return Nickname, nil
	case "preferred_username":
		return PreferredUsername, nil
	case "profile":
		return Profile, nil
	case "picture":
		return Picture, nil
	case "website":
		return Website, nil
	case "gender":
		return Gender, nil
	case "birthdate":
		return Birthdate, nil
	case "zoneinfo":
		return Zoneinfo, nil
	case "locale":
		return Locale, nil
	case "updated_at":
		return UpdatedAt, nil
	case "email":
		return Email, nil
	case "email_verified":
		return EmailVerified, nil
	case "address":
		return Address, nil
	case "phone_number":
		return PhoneNumber, nil
	case "phone_number_verified":
		return PhoneNumberVerified, nil
	default:
		return 0, errors.New("invalid claim type")
	}
}

type Config struct {
	Request  []headerInjector
	Response []headerInjector
}
