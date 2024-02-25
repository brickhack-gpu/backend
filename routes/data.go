package routes

import (
    "context"
    "net/http"

	"github.com/golang-jwt/jwt"

    "gpu/model"
    "gpu/util"
)

type TemplatesRes struct {
	Success   bool             `json:"success"`
	Templates []model.Template `json:"templates"`
}

func (router *Router) Templates(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	templates := []model.Template{}
	err := router.DB.NewSelect().Model(&templates).OrderExpr("active DESC").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

	res := TemplatesRes{
		Templates: templates,
		Success:   true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type SearchRes struct {
	Success   bool             `json:"success"`
	ServerConfigs []model.ServerConfig`json:"server_configs"`
}

func (router *Router) Search(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	serverConfigs := []model.ServerConfig{}
	err := router.DB.NewSelect().Model(&serverConfigs).OrderExpr("active DESC").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

	res := SearchRes{
		ServerConfigs: serverConfigs,
		Success:   true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type DataRes struct {
	Success   bool             `json:"success"`
    Disk int64 `json:"disk"`
    Costs float64 `json:"costs"`
}

func (router *Router) Data(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	uid := int64(props["sub"].(float64))

	diskQuery := router.DB.NewSelect().Model((*model.Product)(nil)).ColumnExpr("SUM(price) AS costs").ColumnExpr("SUM(storage) AS disk").Where("status = 'active'").Where("user_id = ?", uid)
	var costs float64
	var disk int64
	if err := diskQuery.Scan(ctx, &costs, &disk); err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

	res := DataRes{
		Success:   true,
        Disk: disk,
        Costs: costs,
	}

	util.ResJSON(w, http.StatusOK, res)
}
