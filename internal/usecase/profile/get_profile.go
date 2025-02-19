package profile

import (
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
)


type ProfileRes struct {
	Username    string `json:"username"`
	Slug        string `json:"slug"`
	AvatarImage string `json:"avatar_image"`
	CoverImage  string `json:"cover_image"`
	HeadLine    string `json:"head_line"`
}

func GetProfile(r *http.Request, user database.User) (ProfileRes, error) {
	profile, err := apiCfg.DB.GetProfile(r.Context(), user.ID)
	if err != nil {
		return ProfileRes{}, fmt.Errorf("couldn't get profile: %v", err)
	}
	return getProfileRes(profile), nil
}


func getProfileRes(dbProfile database.Profile) ProfileRes {
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
		Username: dbProfile.Username,
		Slug: dbProfile.Slug,
		AvatarImage: imgAvatar,
		CoverImage: imgCover,
		HeadLine: headLine,
	}
}