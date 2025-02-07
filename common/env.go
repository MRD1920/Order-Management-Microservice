package common

import "syscall"

func EnvString(key, fallback string) string {
	if val, found := syscall.Getenv(key); found {
		return val
	}

	return fallback
}
