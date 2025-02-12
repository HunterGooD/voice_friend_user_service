package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func BenchmarkJWT_GenerateAllTokensAsync(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtManager := NewJWTGeneratorWithPrivateKey(key, "tester", 15*time.Second, 30*time.Second, []string{})
	ctx := context.Background()
	uid := "sskdmas"
	role := "user"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtManager.GenerateAllTokensAsync(ctx, uid, role, "")
		if err != nil {
			b.Fatalf("GenerateAllTokensAsync failed: %v", err)
		}
	}
}

func BenchmarkJWT_GenerateAllTokens(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtManager := NewJWTGeneratorWithPrivateKey(key, "tester", 15*time.Second, 30*time.Second, []string{})
	ctx := context.Background()
	uid := "sskdmas"
	role := "user"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtManager.GenerateAllTokens(ctx, uid, role, "")
		if err != nil {
			b.Fatalf("GenerateAllTokens failed: %v", err)
		}
	}
}
