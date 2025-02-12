package auth_test

import (
	"github.com/HunterGooD/voice_friend/user_service/pkg/auth"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	// Инициализируем hasher с параметрами для тестирования
	timeParam := uint32(1)
	memory := uint32(64 * 1024)
	keyLen := uint32(32)
	saltLen := uint32(16)
	threads := uint8(4)

	hasher := auth.NewArgon2Hasher(timeParam, memory, keyLen, saltLen, threads)

	// Исходный пароль для теста
	plainPassword := "supersecret"

	// Генерация хеша
	hashedPassword, err := hasher.HashPassword(plainPassword)
	require.NoError(t, err, "HashPassword не должен возвращать ошибку")
	require.NotEmpty(t, hashedPassword, "Хешированный пароль не должен быть пустым")

	// Дополнительно можно проверить, что строка содержит признаки формата Argon2
	require.True(t, strings.HasPrefix(hashedPassword, "$argon2id$"), "Неверный формат хеша: %s", hashedPassword)

	// Проверяем, что пароль проходит верификацию
	valid, err := hasher.CheckPassword(plainPassword, hashedPassword)
	require.NoError(t, err, "CheckPassword не должен возвращать ошибку для корректного пароля")
	require.True(t, valid, "Пароль должен быть валидным")

	// Проверяем, что неверный пароль не проходит проверку
	valid, err = hasher.CheckPassword("wrongpassword", hashedPassword)
	require.NoError(t, err, "CheckPassword не должен возвращать ошибку для некорректного пароля")
	require.False(t, valid, "Пароль должен быть невалидным")
}

func TestCheckPasswordWithMalformedHash(t *testing.T) {
	// Инициализируем hasher с параметрами для тестирования
	timeParam := uint32(1)
	memory := uint32(64 * 1024)
	keyLen := uint32(32)
	saltLen := uint32(16)
	threads := uint8(4)

	hasher := auth.NewArgon2Hasher(timeParam, memory, keyLen, saltLen, threads)

	// Передаём строку, которая не соответствует ожидаемому формату
	malformedHash := "this-is-not-a-valid-argon2-hash"

	_, err := hasher.CheckPassword("anyPassword", malformedHash)
	require.Error(t, err, "Ожидалась ошибка при передаче некорректного формата хеша")
}
