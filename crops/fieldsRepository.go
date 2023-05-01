package crops

// transaction: https://gist.github.com/miguelmota/d54814683346c4c98cec432cf99506c0

import (
	"context"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type Field struct {
	Id                 uint64 // id
	RegistrationNumber string // registration number in official register
	ownerId            int64  // foreign key
	isVerified         bool
}

type CropType struct {
	Id   uint32 // id
	Name string // crop description/name
}

func insertField(ctx context.Context, field Field) (*Field, error) {
	rows, err := sqldb.Query(ctx, `
		INSERT INTO fields (registrationNumber, ownerId)
		VALUES ($1, $2) RETURNING id
	`, field.RegistrationNumber, field.ownerId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	var id uint64
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	f := Field{Id: id, RegistrationNumber: field.RegistrationNumber, ownerId: field.ownerId, isVerified: false}
	return &f, nil
}

func listFields(ctx context.Context) ([]Field, error) {
	fields := []Field{}
	rows, err := sqldb.Query(ctx, `SELECT id, registrationNumber, ownerId FROM fields where verified == true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := Field{isVerified: true}
		err := rows.Scan(&f.Id, &f.RegistrationNumber, &f.ownerId)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func listFieldsByOwner(ctx context.Context, ownerId int64) ([]Field, error) {
	fields := []Field{}
	rows, err := sqldb.Query(ctx, `SELECT id, registrationNumber, ownerId, verified FROM fields WHERE ownerId = $1`, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := Field{}
		err := rows.Scan(&f.Id, &f.RegistrationNumber, &f.ownerId, &f.isVerified)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func fieldById(ctx context.Context, fieldId uint64) (*Field, error) {
	rows, err := sqldb.Query(ctx, `SELECT id, registrationNumber, ownerId, verified FROM fields WHERE id = $1::int4`, fieldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	f := Field{}
	err = rows.Scan(&f.Id, &f.RegistrationNumber, &f.ownerId, &f.isVerified)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func updateField(ctx context.Context, field Field) (*Field, error) {
	_, err := sqldb.Exec(ctx, `
		UPDATE fields 
		SET registrationNumber = $1, ownerId = $2, verified = $3,
		WHERE id = $4
	`, field.RegistrationNumber, field.ownerId, field.isVerified, field.Id)
	if err != nil {
		rlog.Error("update failed", "err", err)
		return nil, err
	}

	return &field, nil
}

func insertCropType(ctx context.Context, name string) (*CropType, error) {
	rows, err := sqldb.Query(ctx, `
		INSERT INTO cropTypes (name)
		VALUES ($1) RETURNING id
	`, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var id uint32
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}

	ct := CropType{Name: name, Id: id}

	return &ct, nil
}

func listCropTypes(ctx context.Context) ([]CropType, error) {
	cropTypes := []CropType{}
	rows, err := sqldb.Query(ctx, `SELECT id, name FROM cropTypes`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := CropType{}
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			return nil, err
		}
		cropTypes = append(cropTypes, c)
	}
	return cropTypes, nil
}
