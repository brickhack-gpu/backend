package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"gpu/model"
	"gpu/util"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRes struct {
	RefreshToken string `json:"refreshToken"`
	Token        string `json:"token"`
	Success      bool   `json:"success"`
}

func validPassword(password string) bool {
	return len(password) > 6
}

func validUsername(username string) bool {
	return len(username) > 3 && len(username) < 20 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(username)
}

func (router *Router) Register(w http.ResponseWriter, r *http.Request) {
    util.ResError(nil, w, http.StatusBadRequest, "Registration disabled for BrickHack")
    return

	var req RegisterReq
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to decode input.")
		return
	}

	if !validPassword(req.Password) {
		util.ResError(err, w, http.StatusBadRequest, "Invalid password.")
		return
	}

	if !validUsername(req.Username) {
		util.ResError(err, w, http.StatusBadRequest, "Invalid username.")
		return
	}

	exists, err := router.DB.NewSelect().Model((*model.User)(nil)).Where("username = ?", req.Username).Exists(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}
	if exists {
		util.ResError(err, w, http.StatusBadRequest, "Username already in use.")
		return
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Invalid password.")
		return
	}
	user := model.User{
		Username:     req.Username,
		PasswordHash: hashed,
	}
	var id int64
	err = router.DB.NewInsert().Model(&user).Returning("id").Scan(ctx, &id)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to register user.")
		return
	}

	token, refresh, err := util.GenerateJWT(user.Username, id, false, router.JwtSecret)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to generate token.")
		return
	}

	res := RegisterRes{
		Token:        token,
		RefreshToken: refresh,
		Success:      true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	RefreshToken string `json:"refreshToken"`
	Token        string `json:"token"`
	Success      bool   `json:"success"`
}

func (router *Router) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to decode input.")
		return
	}

	u := new(model.User)
	err = router.DB.NewSelect().Model(u).Where("username = ?", req.Username).Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to get user.")
		return
	}

	if !util.CheckPasswordHash(req.Password, u.PasswordHash) {
		util.ResError(err, w, http.StatusBadRequest, "Invalid password.")
		return
	}

	if !u.Active {
		util.ResError(err, w, http.StatusBadRequest, "Disabled.")
		return
	}

	token, refresh, err := util.GenerateJWT(u.Username, u.ID, u.Admin, router.JwtSecret)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to generate token.")
		return
	}

	res := LoginRes{
		Token:        token,
		RefreshToken: refresh,
		Success:      true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type RefreshTokenReq struct {
	Token string `json:"token"`
}

type RefreshTokenRes struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func (router *Router) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenReq

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to decode input.")
		return
	}

	token, err := util.GenerateJWTFromRefreshToken(router.DB, router.JwtSecret, req.Token, context.Background())
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, err.Error())
		return
	}

	res := RefreshTokenRes{
		Token:   token,
		Success: true,
	}

	util.ResJSON(w, http.StatusOK, res)
}
