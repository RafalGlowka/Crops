
CREATE TABLE users(
	id serial PRIMARY KEY,
	email TEXT NOT NULL,
	passwordHash TEXT NOT NULL,
	verifier BOOL NOT NULL,
	seller BOOL NOT NULL,
	UNIQUE(email)
);

CREATE TABLE fields(
	id SERIAL PRIMARY KEY,
	registrationNumber TEXT NOT NULL,
	ownerId INTEGER NOT NULL,
	CONSTRAINT fk_ownerId
		FOREIGN KEY (ownerId) 
		REFERENCES users(id),
	UNIQUE(registrationNumber)
);

CREATE TABLE cropTypes(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	UNIQUE(name)
);

CREATE TABLE offers(
	id SERIAL PRIMARY KEY,
	sellerId INTEGER NOT NULL,
	fieldId INTEGER NOT NULL,
	year INTEGER NOT NULL,
	price BIGINT NOT NULL,
	cropType INTEGER NOT NULL,
	status INTEGER NOT NULL,
	buyerId INTEGER,
	CONSTRAINT fk_sellerId
		FOREIGN KEY (sellerId) 
		REFERENCES users(id),
	CONSTRAINT fk_fieldId
		FOREIGN KEY (fieldId) 
		REFERENCES fields(id),
	CONSTRAINT fk_cropType
		FOREIGN KEY (cropType) 
		REFERENCES cropTypes(id),
	CONSTRAINT fk_buyerId
		FOREIGN KEY (buyerId) 
		REFERENCES users(id),
	UNIQUE(fieldId, year)
);

CREATE TABLE harvests(
	id SERIAL PRIMARY KEY,
	buyerId INTEGER,
	year INTEGER NOT NULL,
	fieldId INT NOT NULL,
	cropType INT NOT NULL,
	price BIGINT NOT NULL,
	harvestSize BIGINT NOT NULL,
	CONSTRAINT fk_buyerId
		FOREIGN KEY (buyerId) 
		REFERENCES users(id),
	CONSTRAINT fk_fieldId
		FOREIGN KEY (fieldId) 
		REFERENCES fields(id),
	CONSTRAINT fk_cropType
		FOREIGN KEY (cropType) 
		REFERENCES cropTypes(id),
	UNIQUE(fieldId, year)		
);