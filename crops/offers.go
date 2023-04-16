package crops

import (
	"context"
	"time"

	"database/sql"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type OffersResponse struct {
	Items []Offer
}

//encore:api auth method=GET path=/offers
func getOffers(ctx context.Context) (*OffersResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	offers, err := listOffers(ctx)
	if err != nil {
		rlog.Error("gettingOffers error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	offersResponse := OffersResponse{Items: offers}
	return &offersResponse, nil
}

//encore:api auth method=GET path=/myOffers
func getMyOffers(ctx context.Context) (*OffersResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	offers, err := listOffersByOwner(ctx, sessionData.UserId)
	if err != nil {
		rlog.Error("getMyOffers error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	offersResponse := OffersResponse{Items: offers}
	return &offersResponse, nil
}

type AddOfferParam struct {
	FieldId  uint64 // field.id
	Year     int
	Price    uint64
	CropType uint32 // CropType.id
	Status   int
}

type OfferResponse struct {
	Item Offer
}

//encore:api auth method=POST path=/offers/add
func addOffer(ctx context.Context, params *AddOfferParam) (*OfferResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	field, err := fieldById(ctx, params.FieldId)
	if err != nil {
		rlog.Error("no field in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "Incorrect field",
		}
	}

	userId := sessionData.UserId
	if field.ownerId != userId {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "Only owner can create an offer",
		}
	}

	currentTime := time.Now()
	if currentTime.Year() > params.Year {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "Only years that in the future",
		}
	}

	f := Offer{FieldId: params.FieldId, SellerId: userId, Year: params.Year, Price: params.Price, CropType: params.CropType, Status: params.Status}
	insertedOffer, err := insertOffer(ctx, f)
	if err != nil {
		rlog.Error("cannot place in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &OfferResponse{Item: *insertedOffer}, nil

}

type BuyOfferParam struct {
	OfferId int64 // offer.id
}

//encore:api auth method=POST path=/offers/buy
func buyOffer(ctx context.Context, params *BuyOfferParam) (*OfferResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	offer, err := offerById(ctx, params.OfferId)
	if err != nil {
		rlog.Error("no offer in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "Incorrect offer",
		}
	}

	userId := sessionData.UserId
	if offer.SellerId == userId {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "You cannot buy your own offer",
		}
	}

	offer.BuyerId = sql.NullInt64{Int64: int64(userId), Valid: true}

	updatedOffer, err := updateOffer(ctx, *offer)
	if err != nil {
		rlog.Error("cannot update in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &OfferResponse{Item: *updatedOffer}, nil

}

type HarvestOfferParam struct {
	OfferId     int64 // offer.id
	HarvestSize uint32
}

type HarvestResponse struct {
	Item Harvest
}

//encore:api auth method=POST path=/offers/harvest
func harvestOffer(ctx context.Context, params *HarvestOfferParam) (*HarvestResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	offer, err := offerById(ctx, params.OfferId)
	if err != nil {
		rlog.Error("no offer in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: "Incorrect offer",
		}
	}

	userId := sessionData.UserId
	if offer.BuyerId.Int64 != userId || offer.BuyerId.Valid == false {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "Only buyer can harvest field",
		}
	}

	harvest := Harvest{BuyerId: offer.BuyerId.Int64, Year: offer.Year, FieldId: offer.FieldId, CropType: offer.CropType, Price: offer.Price, HarvestSize: params.HarvestSize}

	insertedHarvest, err := insertHarvest(ctx, harvest)
	if err != nil {
		rlog.Error("cannot insert into DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &HarvestResponse{Item: *insertedHarvest}, nil

}
