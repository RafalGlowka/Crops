package crops

import (
	"context"

	"database/sql"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type Offer struct {
	Id       uint64 // id
	SellerId int64  // user.id
	FieldId  uint64 // field.id
	Year     int
	Price    uint64
	CropType uint32 // CropType.id
	Status   int
	BuyerId  sql.NullInt64 // user.id
}

func insertOffer(ctx context.Context, offer Offer) (*Offer, error) {
	rows, err := sqldb.Query(ctx, `
		INSERT INTO offers (sellerId, fieldId, year, price, cropType, status, buyerId)
		VALUES ($1, $2, $3, $4, $5, $6, NULL) RETURNING id
	`, offer.SellerId, offer.FieldId, offer.Year, offer.Price, offer.CropType, offer.Status)
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

	o := Offer{Id: id, SellerId: offer.SellerId, FieldId: offer.FieldId, Year: offer.Year, Price: offer.Price, CropType: offer.CropType, Status: offer.Status, BuyerId: offer.BuyerId}
	return &o, nil
}

func offerById(ctx context.Context, offerId int64) (*Offer, error) {
	rows, err := sqldb.Query(ctx, `SELECT id, sellerId, fieldId, year, price, cropType, status, buyerId FROM offers WHERE id = $1::int4`, offerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	o := Offer{}
	err = rows.Scan(&o.Id, &o.SellerId, &o.FieldId, &o.Year, &o.Price, &o.CropType, &o.Status, &o.BuyerId)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func updateOffer(ctx context.Context, offer Offer) (*Offer, error) {
	_, err := sqldb.Exec(ctx, `
		UPDATE offers 
		SET sellerId = $1, fieldId = $2, year = $3, price = $4, cropType = $5, status= $6, buyerId = $7
		WHERE id = $8
	`, offer.SellerId, offer.FieldId, offer.Year, offer.Price, offer.CropType, offer.Status, offer.BuyerId.Int64, offer.Id)
	if err != nil {
		rlog.Error("update failed", "err", err)
		return nil, err
	}

	return &offer, nil
}

func listOffers(ctx context.Context) ([]Offer, error) {
	offers := []Offer{}
	rows, err := sqldb.Query(ctx, `SELECT id, sellerId, fieldId, year, price, cropType, status, buyerId FROM offers WHERE buyerId is NULL`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		o := Offer{}
		err := rows.Scan(&o.Id, &o.SellerId, &o.FieldId, &o.Year, &o.Price, &o.CropType, &o.Status, &o.BuyerId)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}
	return offers, nil

}

func listOffersByOwner(ctx context.Context, ownerId int64) ([]Offer, error) {
	offers := []Offer{}
	rows, err := sqldb.Query(ctx, `SELECT id, sellerId, fieldId, year, price, cropType, status, buyerId FROM offers WHERE sellerId == $1`, ownerId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		o := Offer{}
		err := rows.Scan(&o.Id, &o.SellerId, &o.FieldId, &o.Year, &o.Price, &o.CropType, &o.Status, &o.BuyerId)
		if err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}
	return offers, nil

}
