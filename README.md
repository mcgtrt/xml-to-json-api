# xml-to-json-api
Polling XML feed data &amp; exposing them to the JSON Gateway API

# How to run

1. Clone the repository
```
git clone github.com/mcgtrt/xml-to-json-api
```

2. Go to the project folder
```
cd xml-to-json-api
```

3. Create .env file with your configuration. Example configuration with required files below:
```
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=xmlToJsonApi
HTTP_LISTEN_ADDR=:3000
```

4. Run project
```
make run
```
