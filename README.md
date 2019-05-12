# Supermarket API


## How to run

### Machine Requirements:
* docker 

### Steps To Run
Pull the image from dockerhub. 
Note the latest build of master will be the latest tag.
```
docker pull xmattstrongx/supermarket
```

Run the image choosing a port you would like to expose
```
docker run --rm -it -p 8080:8080 xmattstrongx/supermarket
```

## API Documentation


### GET 

`{{visit_api_host}}/api/states/`

Functionality: Gets a list of all states

### GET 

`{{visit_api_host}}/api/states/{state}/`

Functionality: Get a specific state's details

### GET 

`{{visit_api_host}}/api/states/{state}/cities/`

Functionality: get a specific city's details

Note: I added in an extra large dataset of cities. I would not expose control to the user to add new cities unless required. Instead I grabbed a large dataset from publi records and manipulated it to fit the required cities.csv format. This list should contain every real cities from each state however I have noticed some missing. The extra large dataset was added purposely to show off the pagination funcionality.

### GET 

`{{visit_api_host}}/api/states/AL/cities/{city}`

Functionality: get all cities of a specific state

Optional query string parameters. Choose either of the two combinations. They are all unsigned integers that must be greater than 0
* limit and/or offset
* page and per_page

Note: if both page and per_page are present it will take precendence for pagination else limit and offset will take precendence.

### GET 

`{{visit_api_host}}/api/users/`

Functionality: Gets a list of all users

### GET 

`{{visit_api_host}}/api/users/{userid}/visits/states`

Functionality: Gets all states visited by a specific user

### GET, POST, DELETE

`{{visit_api_host}}/api/users/{userid}/visits`

Functionality: 
* Get all visits for a specific user
* Posts a single new visit for a user. Body must be application/json. The city must exist in the state.
* Deletes a single specific visit for a user

example POST body: 

```
{
    "city": "Carolina beach",
    "state": "NC"
}
```

Optional query string parameters. Choose either of the two combinations. They are all unsigned integers that must be greater than 0
* limit and/or offset
* page and per_page

Note: if both page and per_page are present it will take precendence for pagination else limit and offset will take precendence.

### Key

`{{visit_api_host}}` = The fully URL and port where the api is hosted. Example `http://dockerhost:8080`

`{state}` = The abbreviation of the state name. Case insensitive. Example `NC`

`{city}` = The full city name. Case insensitive. Example `Carolina Beach`

`{visitid}` = The unique ID of the visit. GUID. Example = `ad421bd4-cbcb-11e6-a819-0242ac140003`

`{userid}` = The unique ID of the user. GUID. Example = `9e77a0fc-cbbb-11e6-8d18-0242ac140003`
