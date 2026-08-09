package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/identicalaffiliation/xakat0n/backend/internal/modules/auth/domain"
)

type fakeUserRepository struct {
	userID uuid.UUID
	err    error
	calls  []domain.Username
}

func (f *fakeUserRepository) GetOrCreate(ctx context.Context, username domain.Username) (uuid.UUID, error) {
	f.calls = append(f.calls, username)
	if f.err != nil {
		return uuid.Nil, f.err
	}

	return f.userID, nil
}

type fakeTokenIssuer struct {
	token string
	err   error
}

func (f *fakeTokenIssuer) Issue(userID uuid.UUID, username domain.Username) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	return f.token, nil
}

type fakeLogger struct{}

func (l *fakeLogger) DebugContext(context.Context, string, ...any)                   {}
func (l *fakeLogger) InfoContext(context.Context, string, ...any)                    {}
func (l *fakeLogger) WarnContext(context.Context, string, ...any)                    {}
func (l *fakeLogger) ErrorContext(context.Context, string, ...any)                   {}
func (l *fakeLogger) WithField(ctx context.Context, _ string, _ any) context.Context { return ctx }
func (l *fakeLogger) WrapError(_ context.Context, err error) error                   { return err }
func (l *fakeLogger) ContextFromError(ctx context.Context, _ error) context.Context  { return ctx }

func TestLoginUsecase(t *testing.T) {
	t.Parallel()

	t.Run("success normalizes username before lookup and issuing", func(t *testing.T) {
		userID := uuid.New()
		users := &fakeUserRepository{userID: userID}
		issuer := &fakeTokenIssuer{token: "jwt-token"}
		usecase := NewLoginUsecase(users, issuer, &fakeLogger{})

		response, err := usecase.Login(context.Background(), "  Hunter_42  ")
		require.NoError(t, err)
		require.Len(t, users.calls, 1)
		assert.Equal(t, "hunter_42", users.calls[0].String())
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, "hunter_42", response.Username)
		assert.Equal(t, "jwt-token", response.Token)
	})

	t.Run("invalid username never reaches the repository", func(t *testing.T) {
		users := &fakeUserRepository{}
		usecase := NewLoginUsecase(users, &fakeTokenIssuer{}, &fakeLogger{})

		_, err := usecase.Login(context.Background(), "ab")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrUsernameLength)
		assert.Empty(t, users.calls)
	})

	t.Run("repository error is wrapped as internal", func(t *testing.T) {
		users := &fakeUserRepository{err: errors.New("db down")}
		usecase := NewLoginUsecase(users, &fakeTokenIssuer{}, &fakeLogger{})

		_, err := usecase.Login(context.Background(), "hunter_42")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})

	t.Run("issuer error is wrapped as internal", func(t *testing.T) {
		users := &fakeUserRepository{userID: uuid.New()}
		issuer := &fakeTokenIssuer{err: errors.New("signing failed")}
		usecase := NewLoginUsecase(users, issuer, &fakeLogger{})

		_, err := usecase.Login(context.Background(), "hunter_42")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInternal)
	})
}
