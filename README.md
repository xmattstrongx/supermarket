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

The API is documented using swagger openapi spec 3.0.

To view the swagger api documentation run the swagger make target
```
make swagger
docker run --rm -p 80:8080 -e "SWAGGER_JSON=/spec/swagger.yaml" -v /Users/matthewstrong/go/modules/github.com/xmattstrongx/supermarket/swagger/spec:/spec swaggerapi/swagger-ui
```

Next open a browser to localhost

![Swagger Example](images/swagger_example.png)

## Devflow

To test, build and run the code locally use the deploy make target.

```
make deploy
```