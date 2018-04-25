package authInterface

type AuthInterface interface{
	Valid(l, p string) bool
}
