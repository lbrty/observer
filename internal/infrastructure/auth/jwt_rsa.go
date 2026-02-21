package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

// Claims are the JWT claims used in this application.
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Type   string `json:"type"`
}

// TokenGenerator defines the token operations interface.
type TokenGenerator interface {
	GenerateAccessToken(userID ulid.ULID, role string) (string, time.Time, error)
	GenerateRefreshToken() (string, error)
	GenerateMFAToken(userID ulid.ULID) (string, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateMFAToken(tokenString string) (*Claims, error)
}

// RSATokenGenerator implements TokenGenerator using RSA keys.
type RSATokenGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessTTL  time.Duration
	refreshTTL time.Duration
	mfaTTL     time.Duration
	issuer     string
}

// NewRSATokenGenerator creates a new RSATokenGenerator.
func NewRSATokenGenerator(keys *RSAKeys, accessTTL, refreshTTL, mfaTTL time.Duration, issuer string) *RSATokenGenerator {
	return &RSATokenGenerator{
		privateKey: keys.PrivateKey,
		publicKey:  keys.PublicKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		mfaTTL:     mfaTTL,
		issuer:     issuer,
	}
}

// GenerateAccessToken creates a signed access JWT for the given user.
func (g *RSATokenGenerator) GenerateAccessToken(userID ulid.ULID, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(g.accessTTL)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: userID.String(),
		Role:   role,
		Type:   "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return tokenStr, expiresAt, nil
}

// GenerateRefreshToken returns a secure random ULID string as the refresh token.
func (g *RSATokenGenerator) GenerateRefreshToken() (string, error) {
	return ulid.Make().String(), nil
}

// GenerateMFAToken creates a short-lived JWT for MFA verification.
func (g *RSATokenGenerator) GenerateMFAToken(userID ulid.ULID) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.mfaTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID: userID.String(),
		Type:   "mfa_pending",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.privateKey)
}

// ValidateAccessToken parses and validates an access token.
func (g *RSATokenGenerator) ValidateAccessToken(tokenString string) (*Claims, error) {
	return g.validateToken(tokenString, "access")
}

// ValidateMFAToken parses and validates an MFA pending token.
func (g *RSATokenGenerator) ValidateMFAToken(tokenString string) (*Claims, error) {
	return g.validateToken(tokenString, "mfa_pending")
}

func (g *RSATokenGenerator) validateToken(tokenString, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
	}

	return claims, nil
}
