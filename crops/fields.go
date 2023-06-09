package crops

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type FieldsResponse struct {
	Items []Field
}

//encore:api auth method=GET path=/fields
func getFields(ctx context.Context) (*FieldsResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	fields, err := listFields(ctx)
	if err != nil {
		rlog.Error("gettingFields error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	fieldsResponse := FieldsResponse{Items: fields}
	return &fieldsResponse, nil
}

//encore:api auth method=GET path=/myFields
func getMyFields(ctx context.Context) (*FieldsResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	fields, err := listFieldsByOwner(ctx, sessionData.UserId)
	if err != nil {
		rlog.Error("gettingMyFields error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	fieldsResponse := FieldsResponse{Items: fields}
	return &fieldsResponse, nil
}

type AddFieldParam struct {
	RegistrationNumber string
}

type FieldResponse struct {
	Item Field
}

//encore:api auth method=POST path=/fields/add
func addField(ctx context.Context, params *AddFieldParam) (*FieldResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	registrationNumberValidationResult := createRegistrationNumberValidator().validate(params.RegistrationNumber)
	if registrationNumberValidationResult != nil {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: registrationNumberValidationResult.errorMessage,
		}
	}

	userId := sessionData.UserId
	f := Field{RegistrationNumber: params.RegistrationNumber, ownerId: userId}
	insertedField, err := insertField(ctx, f)
	if err != nil {
		rlog.Error("cannot place in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &FieldResponse{Item: *insertedField}, nil

}

type VerifyFieldParam struct {
	FieldId uint64
}

//encore:api auth method=POST path=/fields/verify
func verifyField(ctx context.Context, params *VerifyFieldParam) (*FieldResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	userId := sessionData.UserId
	user, err := getUserById(ctx, userId)
	if err != nil || user.isVerifier == false {
		rlog.Error("cannot place in DB", "err", err)
		var message = ""
		if err != nil {
			message = "Missing permissions"
		} else {
			message = err.Error()
		}
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: message,
		}
	}

	field, err := fieldById(ctx, params.FieldId)
	if err != nil {
		rlog.Error("cannot find field in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: err.Error(),
		}
	}

	field.isVerified = true

	insertedField, err := updateField(ctx, *field)
	if err != nil {
		rlog.Error("cannot place in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &FieldResponse{Item: *insertedField}, nil

}

type AddCropTypeParam struct {
	Name string
}

type CropTypeResponse struct {
	Item CropType
}

//encore:api auth method=POST path=/crops/add
func addCropType(ctx context.Context, params *AddCropTypeParam) (*CropTypeResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	userId := sessionData.UserId
	user, err := getUserById(ctx, userId)
	if err != nil || user.isVerifier == false {
		rlog.Error("cannot place in DB", "err", err)
		var message = ""
		if err != nil {
			message = "Missing permissions"
		} else {
			message = err.Error()
		}
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: message,
		}
	}

	registrationNumberValidationResult := createRegistrationNumberValidator().validate(params.Name)
	if registrationNumberValidationResult != nil {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: registrationNumberValidationResult.errorMessage,
		}
	}

	insertedCropType, err := insertCropType(ctx, params.Name)
	if err != nil {
		rlog.Error("cannot place in DB", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}

	return &CropTypeResponse{Item: *insertedCropType}, nil

}

type FieldHistoryParam struct {
	FieldId uint64
}

type FieldHistoryResponse struct {
	Items []Harvest
}

//encore:api auth method=POST path=/fields/history
func fieldHistory(ctx context.Context, params *FieldHistoryParam) (*FieldHistoryResponse, error) {
	uid, authenticated := auth.UserID()
	sessionData := auth.Data().(*SessionData)

	if len(uid) < 1 || authenticated != true || sessionData == nil || sessionData.UserId == 0 {
		rlog.Error("not authenticated", "authenticatedFlag", authenticated)
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "Permission denied",
		}
	}

	harvests, err := listHarvest(ctx)
	if err != nil {
		rlog.Error("gettingHarvest error!", "err", err)
		return nil, &errs.Error{
			Code:    errs.Aborted,
			Message: err.Error(),
		}
	}
	fieldsResponse := FieldHistoryResponse{Items: harvests}
	return &fieldsResponse, nil

}
