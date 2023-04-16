package crops

import (
	"context"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type Harvest struct {
	Id          uint64 // id
	BuyerId     int64  // user.id
	Year        int
	FieldId     uint64 // field.id
	CropType    uint32 //cropType.id
	Price       uint64
	HarvestSize uint32
}

func insertHarvest(ctx context.Context, harvest Harvest) (*Harvest, error) {
	rows, err := sqldb.Query(ctx, `
		INSERT INTO harvests (buyerId, year, fieldId, cropType, price, harvestSize)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, harvest.BuyerId, harvest.Year, harvest.FieldId, harvest.CropType, harvest.Price, harvest.HarvestSize)
	if err != nil {
		rlog.Error("insert failed", "err", err)
		return nil, err
	}

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		rlog.Error("scan failed", "err", err)
		return nil, err
	}

	rows.Next()
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		rlog.Error("scan failed", "err", err)
		return nil, err
	}

	h := Harvest{Id: id, BuyerId: harvest.BuyerId, FieldId: harvest.FieldId, Year: harvest.Year, Price: harvest.Price, CropType: harvest.CropType, HarvestSize: harvest.HarvestSize}
	return &h, nil
}

func listHarvest(ctx context.Context) ([]Harvest, error) {
	harvests := []Harvest{}
	rows, err := sqldb.Query(ctx, `SELECT id, buyerId, year, fieldId, cropType, price, harvestSize FROM harvests`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		h := Harvest{}
		err := rows.Scan(&h.Id, &h.BuyerId, &h.Year, &h.FieldId, &h.CropType, &h.Price, &h.HarvestSize)
		if err != nil {
			return nil, err
		}
		harvests = append(harvests, h)
	}
	return harvests, nil
}
