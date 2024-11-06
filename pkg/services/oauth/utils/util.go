package utils

import "fmt"

func GetCallbackUrl(domain string) string {
	return fmt.Sprintf("%s%s", domain, "api/callback")
}
