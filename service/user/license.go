package user

import "time"

func (u *UserService) CheckLicenseValid(user *User) bool {
	if user.AccountType == TestAccount {
		return true
	}
	if time.Now().After(user.LicenseExpiry) {
		return false
	}
	return true
}
