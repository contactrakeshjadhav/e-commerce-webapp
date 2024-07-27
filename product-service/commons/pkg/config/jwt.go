package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	jwtSigningKey       = "JWT_SIGNING_KEY"
	jwtTokenKey         = "JWT_TOKEN_KEY"
	jwtExpiryTime       = "JWT_EXPIRY_TIME"
	serviceAccountToken = "SERVICE_ACCOUNT_TOKEN"
)

type JWTConfig struct {
	SigningKey          string
	TokenKey            string
	ServiceAccountToken string
	ExpiryTime          int
}

func LoadJWTConfig() (JWTConfig, error) {
	signingKey := os.Getenv(jwtSigningKey)
	if signingKey == "" {
		return JWTConfig{}, errors.New("failed to load signing key")
	}

	tokenKey := os.Getenv(jwtTokenKey)
	if tokenKey == "" {
		return JWTConfig{}, errors.New("failed to load token key")
	}

	expiryTimeValue := os.Getenv(jwtExpiryTime)
	if expiryTimeValue == "" {
		return JWTConfig{}, errors.New("failed to load expiry")
	}

	serviceAccount := os.Getenv(serviceAccountToken)
	if serviceAccount == "" {
		return JWTConfig{}, errors.New("failed to load service account token")
	}

	expiryTime, err := strconv.Atoi(expiryTimeValue)
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{
		SigningKey:          signingKey,
		TokenKey:            tokenKey,
		ExpiryTime:          expiryTime,
		ServiceAccountToken: serviceAccount,
	}, nil
}
