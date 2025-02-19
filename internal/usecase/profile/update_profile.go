package profile

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/fileManager"
)

type UpdateProfileReq struct {
	Slug        string  `json:"slug"`
	Username    string  `json:"username"`
	HeadLine    string  `json:"head_line"`
	AvatarImage *string `json:"avatar_image"`
	CoverImage  *string `json:"cover_image"`
}

func UpdateProfile(r *http.Request, params UpdateProfileReq, user database.User) (ProfileRes, error) {
	profileData, _ := apiCfg.DB.GetProfile(r.Context(), user.ID)
	parameters := database.UpdateProfileParams{
		UserID:   user.ID,
		Slug:     params.Slug,
		Username: params.Username,
		HeadLine: sql.NullString{
			Valid:  true,
			String: params.HeadLine,
		},
		AvatarImage: sql.NullString{
			Valid:  true,
			String: profileData.AvatarImage.String,
		},
		CoverImage: sql.NullString{
			Valid:  true,
			String: profileData.CoverImage.String,
		},
	}
	if params.AvatarImage != nil {
		link, err := fileManager.UploadFile(user.ID, *params.AvatarImage)
		if err != nil {
			return ProfileRes{}, fmt.Errorf("couldn't upload avatar image: %v", err)
		}
		parameters.AvatarImage = sql.NullString{
			Valid:  true,
			String: link,
		}
	}
	if params.CoverImage != nil {
		link, err := fileManager.UploadFile(user.ID, *params.CoverImage)
		if err != nil {
			return ProfileRes{}, fmt.Errorf("couldn't upload cover image: %v", err)
		}
		parameters.CoverImage = sql.NullString{
			Valid:  true,
			String: link,
		}
	}

	profile, err := apiCfg.DB.UpdateProfile(r.Context(), parameters)
	if err != nil {
		return ProfileRes{}, fmt.Errorf("couldn't update profile: %v", err)
	}
	return getProfileRes(profile), nil
}