package routes

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"

	"gpu/model"
	"gpu/util"
)

type ProfileRes struct {
	Success      bool   `json:"success"`
    Balance float64 `json:"balance"`
    Notifications []model.Notification `json:"notifications"`
}

func (router *Router) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
    props, _ := r.Context().Value("props").(jwt.MapClaims)
    uid := int(props["sub"].(float64))

    notifications := []model.Notification{}
	err := router.DB.NewSelect().Model(&notifications).Where("user_id = ?", uid).Where("read = false").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

    depositQuery := router.DB.NewSelect().Model((*model.Deposit)(nil)).ColumnExpr("SUM(amount) AS deposit_sum").Where("user_id = ?", uid)
    var depositSum float64
    if err := depositQuery.Scan(ctx, &depositSum); err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
    }
    purchaseQuery := router.DB.NewSelect().Model((*model.Purchase)(nil)).ColumnExpr("SUM(amount) AS purchase_sum").Where("user_id = ?", uid)
    var purchaseSum float64
    if err := purchaseQuery.Scan(ctx, &purchaseSum); err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
    }

	res := ProfileRes{
        Balance: depositSum - purchaseSum,
        Notifications: notifications,
		Success:      true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type TransactionsRes struct {
	Success      bool   `json:"success"`
    Deposits []*model.Deposit `json:"deposits"`
    Purchases []*model.Purchase `json:"purchases"`
}

func (router *Router) Transactions(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
    props, _ := r.Context().Value("props").(jwt.MapClaims)
    uid := int64(props["sub"].(float64))

    user := model.User{
        ID: uid,
    }
    err := router.DB.NewSelect().Model(&user).WherePK().Relation("Purchases").Relation("Deposits").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

    res := TransactionsRes{
        Deposits: user.Deposits,
        Purchases: user.Purchases,
        Success: true,
    }

    util.ResJSON(w, http.StatusOK, res)
}

type ProductsRes struct {
	Success      bool   `json:"success"`
    Products []*model.Product `json:"products"`
}

func (router *Router) Products(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
    props, _ := r.Context().Value("props").(jwt.MapClaims)
    uid := int64(props["sub"].(float64))

    user := model.User{
        ID: uid,
    }
    err := router.DB.NewSelect().Model(&user).WherePK().Relation("Products").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

    res := ProductsRes{
        Products: user.Products,
        Success: true,
    }

    util.ResJSON(w, http.StatusOK, res)
}

type TemplatesRes struct {
	Success      bool   `json:"success"`
    Templates []model.Template `json:"templates"`
}

func (router *Router) Templates(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

    templates := []model.Template{}
    err := router.DB.NewSelect().Model(&templates).Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

    res := TemplatesRes{
        Templates: templates,
        Success: true,
    }

    util.ResJSON(w, http.StatusOK, res)
}
