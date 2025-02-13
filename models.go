package main

import (
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/google/uuid"
)

type UserRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Email      string `json:"email"`
}

type ProfileRes struct {
	ID          uuid.UUID `json:"profile_uuid"`
	Username    string `json:"username"`
	Slug        string `json:"slug"`
	AvatarImage string `json:"avatar_image"`
	CoverImage  string `json:"cover_image"`
	HeadLine    string `json:"head_line"`
}

func databaseUserToUser(dbUser database.User) UserRes {
	return UserRes{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
}

func databaseProfileToProfile(dbProfile database.Profile) ProfileRes {
	imgAvatar := ""
	imgCover := ""
	headLine := ""
	if dbProfile.AvatarImage.Valid {
		imgAvatar = dbProfile.AvatarImage.String
	}
	if dbProfile.CoverImage.Valid {
		imgCover = dbProfile.CoverImage.String
	}
	if dbProfile.HeadLine.Valid {
		headLine = dbProfile.HeadLine.String
	}
	return ProfileRes{
		ID: dbProfile.ID,
		Username: dbProfile.Username,
		Slug: dbProfile.Slug,
		AvatarImage: imgAvatar,
		CoverImage: imgCover,
		HeadLine: headLine,
	}
}

type LoginUserRes struct{
	SessionID string `json:"session_id"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User UserRes `json:"user"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh"`
}

type RenewAccessTokenRes struct {
	AccessToken string `json:"access"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
