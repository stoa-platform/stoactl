// SPDX-License-Identifier: Apache-2.0
// Copyright 2024-2026 CAB Ingénierie / Christophe ABOULICAM
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/stoa-platform/stoactl/internal/config"
	"github.com/stoa-platform/stoactl/internal/output"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display authentication status",
		Long: `Display the current authentication status.

Shows whether you are authenticated, token expiration,
and basic user information.

Example:
  stoactl auth status`,
		Args: cobra.NoArgs,
		RunE: runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, err := cfg.GetCurrentContext()
	if err != nil {
		return err
	}

	tokenCache, err := config.LoadTokenCache()
	if err != nil {
		return fmt.Errorf("failed to load token cache: %w", err)
	}

	if tokenCache == nil || tokenCache.AccessToken == "" {
		output.Info("Not authenticated.")
		output.Info("Run 'stoactl auth login' to authenticate.")
		return nil
	}

	// Check if token is for current context
	if tokenCache.Context != ctx.Name {
		output.Info("Token is for a different context (%s).", tokenCache.Context)
		output.Info("Run 'stoactl auth login' to authenticate to %s.", ctx.Name)
		return nil
	}

	// Check if token is expired
	if time.Now().Unix() > tokenCache.ExpiresAt {
		output.Info("Token expired.")
		output.Info("Run 'stoactl auth login' to re-authenticate.")
		return nil
	}

	// Parse token claims
	claims, err := parseJWTClaims(tokenCache.AccessToken)
	if err != nil {
		output.Info("Authenticated (unable to parse token details)")
		return nil
	}

	fmt.Println("Authentication Status:")
	fmt.Printf("  Context:    %s\n", ctx.Name)
	fmt.Printf("  Server:     %s\n", ctx.Context.Server)
	fmt.Printf("  Tenant:     %s\n", ctx.Context.Tenant)
	fmt.Println()
	fmt.Println("User:")
	if claims.Email != "" {
		fmt.Printf("  Email:      %s\n", claims.Email)
	}
	if claims.Name != "" {
		fmt.Printf("  Name:       %s\n", claims.Name)
	}
	fmt.Printf("  Subject:    %s\n", claims.Subject)
	fmt.Println()
	fmt.Println("Token:")
	fmt.Printf("  Expires:    %s\n", time.Unix(tokenCache.ExpiresAt, 0).Format(time.RFC3339))
	fmt.Printf("  Valid for:  %s\n", time.Until(time.Unix(tokenCache.ExpiresAt, 0)).Round(time.Minute))

	return nil
}

func parseJWTClaims(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload := parts[1]
	// Add padding if needed
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out from STOA Platform",
		Long: `Remove stored authentication credentials.

Example:
  stoactl auth logout`,
		Args: cobra.NoArgs,
		RunE: runLogout,
	}
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.SaveTokenCache(&config.TokenCache{}); err != nil {
		return fmt.Errorf("failed to clear token cache: %w", err)
	}

	output.Success("Logged out successfully.")
	return nil
}
