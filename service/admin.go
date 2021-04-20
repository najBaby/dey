package service

import "net/http"

// AdminService is
type AdminService uint8

const (
	// AdminSignUp is
	AdminSignUp AdminService = iota
	// AdminSignIn is
	AdminSignIn
	// AdminSignOut is
	AdminSignOut
)

// AdminStatus is
type AdminStatus uint8

const (
	// AdminSuccess is
	AdminSuccess AdminStatus = iota
	// AdminFailure is
	AdminFailure
	// AdminForbidden is
	AdminForbidden
)

// AdminServiceServer is
type AdminServiceServer struct {
}

// SignIn is
func (AdminServiceServer) SignIn(w http.ResponseWriter, r *http.Request) {
}

// SignOut is
func (AdminServiceServer) SignOut(w http.ResponseWriter, r *http.Request) {
}
