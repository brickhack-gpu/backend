package scan

import (
    "log"
    "time"
    "context"

	"github.com/uptrace/bun"

    "gpu/model"
)

func ScanBalance(db *bun.DB) error {
    ctx := context.Background()
    for {
        var activeProducts []model.Product
        err := db.NewSelect().Model(&activeProducts).Where("status = 'active'").Scan(ctx)
        if err != nil {
            log.Println(err)
            time.Sleep(time.Minute * 1)
            continue
        }

        for _, product := range activeProducts {
            purchase := model.Purchase{
                UserID: product.UserID,
                ProductID: product.ID,
                Amount: product.Price,
            }
            _, err := db.NewInsert().Model(&purchase).Exec(ctx)
            if err != nil {
                log.Println(err)
            }
        }

        time.Sleep(time.Minute * 60)
    }
}
