package routes

import (
    "log"
    "fmt"
    "time"
	"context"
	"encoding/json"
	"net/http"

    "github.com/goombaio/namegenerator"
	"github.com/golang-jwt/jwt"

	"gpu/model"
	"gpu/util"
)

type SpinServerReq struct {
    ServerConfigID int64 `json:"server_config_id"`
    TemplateID int64 `json:"template_id"`
    Storage int `json:"storage"`
}

type SpinServerRes struct {
	Success      bool   `json:"success"`
}

func (router *Router) SpinServer(w http.ResponseWriter, r *http.Request) {
	var req SpinServerReq
	ctx := context.Background()
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	uid := int64(props["sub"].(float64))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to decode input.")
		return
	}

    if req.Storage < 100 {
		util.ResError(err, w, http.StatusBadRequest, "This template needs 100GB.")
		return
    }

    if uid != 1 {
		util.ResError(err, w, http.StatusBadRequest, "Insufficent balance.")
		return
    }

    template := model.Template{
        ID: req.TemplateID,
    }
    err = router.DB.NewSelect().Model(&template).WherePK().Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Invalid template.")
		return
	}

    serverConfig := model.ServerConfig{
        ID: req.ServerConfigID,
    }
    err = router.DB.NewSelect().Model(&serverConfig).WherePK().Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Invalid config.")
		return
	}

    // func CreateInstance(projectID, zone, instanceName, machineType, sourceImage, region, script, gpuType string, gpuCount int32, disk int64) error {
    seed := time.Now().UTC().UnixNano()
    nameGenerator := namegenerator.NewNameGenerator(seed)
    gcpId := nameGenerator.Generate()

    product := model.Product{
        Price: serverConfig.Price,
        Status: "spinning",
        GCPID: gcpId,
        Storage: req.Storage,
        UserID: uid,
        ServerConfigID: serverConfig.ID,
        TemplateID: template.ID,
    }
    _, err = router.DB.NewInsert().Model(&product).Exec(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
	}

    go func() {
        err := util.CreateInstance("siggpu", serverConfig.Zone, gcpId, serverConfig.MachineType, template.Container, serverConfig.Region, "",
            serverConfig.GPUType, int32(serverConfig.GPUCount), int64(req.Storage))
        if err != nil {
            product.Status = "failed"
            return
        } else {
            product.Status = "building"
            ip, err := util.GetInstanceIP("siggpu", serverConfig.Zone, gcpId)
            if err != nil {
                log.Println(err)
            }
            product.DNSLink = fmt.Sprintf("http://%s", ip)
        }

        _, err = router.DB.NewUpdate().Model(&product).Where("gcp_id = ?", gcpId).Exec(ctx)
        if err != nil {
            panic(err)
        }

        if product.Status == "building" {
            time.Sleep(10)
            product.Status = "active"
            _, err = router.DB.NewUpdate().Model(&product).Where("gcp_id = ?", gcpId).Exec(ctx)
            if err != nil {
                panic(err)
            }
        }
    }()

    purchase := model.Purchase{
        UserID: product.UserID,
        ProductID: product.ID,
        Amount: product.Price,
    }
    _, err = router.DB.NewInsert().Model(&purchase).Exec(ctx)
    if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Database error.")
		return
    }

	res := SpinServerRes{
		Success:      true,
	}

	util.ResJSON(w, http.StatusOK, res)
}

type KillServerReq struct {
    GCPId string `json:"gcp_id"`
}

func (router *Router) KillServer(w http.ResponseWriter, r *http.Request) {
	var req KillServerReq
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Failed to decode input.")
		return
	}

    product := new(model.Product)
    err = router.DB.NewSelect().Model(product).Where("gcp_id = ?", req.GCPId).Relation("ServerConfig").Scan(ctx)
	if err != nil {
		util.ResError(err, w, http.StatusBadRequest, "Invalid config.")
		return
	}

    go func() {
        err := util.DeleteInstance("siggpu", product.ServerConfig.Zone, req.GCPId)
        if err != nil {
            product.Status = "failed"
            return
        } else {
            product.Status = "destroyed"
        }

        _, err = router.DB.NewUpdate().Model(product).Where("gcp_id = ?", req.GCPId).Exec(ctx)
        if err != nil {
            panic(err)
        }
    }()

    product.Status = "destroying"
    _, err = router.DB.NewUpdate().Model(product).Where("gcp_id = ?", req.GCPId).Exec(ctx)
    if err != nil {
        panic(err)
    }

	res := SpinServerRes{
		Success:      true,
	}

	util.ResJSON(w, http.StatusOK, res)
}
