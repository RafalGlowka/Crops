# Crops Market

Example backend build in Go Lang with the use of encore library. 

It's a simple POC to test few aspects in practice:
- Simple user management with email verification, sessions and balance withing the system.
- Simple record of fields (piece of land) owned by users.
- Each field need to be verified by special user with feature "verifier".
- Simple market with list of crops on fields, that other users may buy. 
- Register of harvested crops - Historical data allows the purchaser to assess the value of the field.

 



## Running

<img width="200px" src="https://encore.dev/assets/branding/logo/logo.svg" alt="Encore - The Backend Development Engine" />

To run the application locally, make sure you have [Encore](https://encore.dev) and Docker](https://docker.com) installed and running. This is required to run Encore applications with SQL databases.

```bash
# Run the app
encore run
```

## Connecting to Database

```bash
encore db conn-uri crops
```

## Using the API

Basic set of postman requests is delivered in "Crop.postman_collection.json" file.
Endpoints with public adnotation do not need session. Endpoints with auth adnotation require header "Authorization" with value "Bearer  <token>". Token is returned in login response, in field "Token".

### User Management

#### public POST http://localhost:4000/user/create
Params: 
	Email - 
	Password - Password should have at least 5 characters
Response:

#### public GET http://localhost:4000/user/verify/:code
Endpoint used for email verification. It's showing the simplest webpage with verification status. 


#### public POST http://localhost:4000/user/login
Params: 
	Email - 
	Password - Password should have at least 5 characters
Response:`
	Token - Session token. Add prefix "Bearer " and send in auth header as "Authorization" field.
	Email - 
	IsVerified - 
	IsVerifier - 
	
#### auth POST http://localhost:4000/user/topUp
Params:
	Amount - Amount in cents that will be added to your balance. It's a simulation of payment methods. 
Response:
	BalanceAfter -

Make user a verifier can be done only by direct connection to Database.

### Crops and fields

auth POST http://localhost:4000/crops/add - Only Verifier can add crops
Params:
	Name - Name of crop.
Response:
	Item -
		Id - 
		Name - 
	
auth GET http://localhost:4000/fields
Params: None
Response:
	Items: [
		Id                 -
		RegistrationNumber -
	]

Only verified fields are listed here 

auth GET http://localhost:4000/myFields
Params: None
Response:

auth POST http://localhost:4000/fields/add
Params: 
	RegistrationNumber - Registration number in country register. 
Response:
	Item: {
		Id                 -
		RegistrationNumber - 
	}

auth POST http://localhost:4000/fields/verify - Only Verifier can verify the field.
Params:
	FieldId - 
Response:
	Item {
		Id                 - 
		RegistrationNumber - 
		ownerId            - 
		isVerified         -
	}

### Offers

auth GET http://localhost:4000/offers
	Only verified fields are visible for other users. 
Params:
Response:
	Items[
		{
			Id       - 
			SellerId -
			FieldId  -
			Year     -
			Price    -
			CropType -
			Status   -
			BuyerId  -
		}
	]
auth GET http://localhost:4000/myOffers
	List of all your fields - verified, sold ect. 
Params:
Response:
	Items[
		{
			Id       - 
			SellerId -
			FieldId  -
			Year     -
			Price    -
			CropType -
			Status   -
			BuyerId  -
		}
	]
auth POST http://localhost:4000/offers/add
	Adding offer for your fields and crops on it. 
Params:
	FieldId  -
	Year     -
	Price    -
	CropType -
	Status   -
Response:
	Item {
		Id       - 
		SellerId -
		FieldId  -
		Year     -
		Price    -
		CropType -
		Status   -
		BuyerId  -
	}
}
auth POST http://localhost:4000/offers/buy
	Buying offer. 
Params:
	OfferId -
Response:
	Item {
		Id       - 
		SellerId -
		FieldId  -
		Year     -
		Price    -
		CropType -
		Status   -
		BuyerId  -
	}
auth POST http://localhost:4000/offers/harvest
	Harvest crops from field that you have bought. 
Params:
	OfferId     -
	HarvestSize -
Response:
	Item {
		Id          -
		BuyerId     -
		Year        -
		FieldId     -
		CropType    -
		Price       -
		HarvestSize -
	}

auth POST http://localhost:4000/fields/history
	Showing history of harvests 
Params:
	OfferId     -
Response:
	Items[
		 {
			Id          -
			BuyerId     -
			Year        -
			FieldId     -
			CropType    -
			Price       -
			HarvestSize -
		}
	]
## Open the developer dashboard

While `encore run` is running, open <http://localhost:4000/> to view Encore's local developer dashboard.
Here you can see the request you just made and a view a trace of the response.

## Deployment

Deploy your application to a staging environment in Encore's free development cloud.

```bash
git push encore
```

Then head over to <https://app.encore.dev> to find out your production URL, and off you go into the clouds!

## Testing

```bash
encore test -v
```

## Future developement
* Login cooldown to prevent bruteforce password attack
* Logout
* Cover with tests: fields and offers