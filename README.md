# gu-port-backend
A REST backend for https://github.com/Flaque/gu-port

## To set up...
Install [golang](https://golang.org/doc/install).

Run the following commands:
```bash
git clone https://github.com/maxchehab/gu-port-backend.git
cd gu-port-backend
go get && ./build
```

Or if you already have go setup, you can run 
``` 
go get github.com/maxchehab/gu-port-backend
```

Access the api at localhost:8080

## API

### Pagination `GET`
#### Request:
`localhost:8080/pages/{count}/{offset}`
#### Response:
```json
{
   "pages":[
      {
         "name":"Another Title",
         "body":"# This is a body",
         "author":"Rick Sanchez",
         "pageID":"1",
         "userID":"0"
      }
   ],
   "count":1,
   "offset":1,
   "total":3
}
```
### Page `GET`
#### Request:
`localhost:8080/users/{userID}/pages/{pageID}`
#### Response:
```json
{
   "name":"Another Title",
   "body":"# This is a body",
   "author":"Rick Sanchez",
   "pageID":"1",
   "userID":"0"
}
```

### Register `POST`
#### Request:
`localhost:8080/register/`
Parameters:
```json
access_code: test
email: maxchehab@gmail.com
password: mypassword
username: maxchehab
```
Body:
`email=maxchehab%40gmail.com&username=mchehab&access_code=test&password=mypassword`
#### Response:
```json
{
   "session":"828f9691-2f4b-4d3d-93e0-3494d55944af",
   "valid":true
}
```
```json
{
   "username":{
      "valid":false,
      "message":[
         "Username taken.",
         "Username cannot be empty."
      ]
   },
   "email":{
      "valid":false,
      "message":[
         "Email taken.",
         "Invalid email."
      ]
   },
   "password":{
      "valid":true,
      "message":[
         "Password cannot be empty."
      ]
   },
   "accessCode":{
      "valid":false,
      "message":[
         "Access code invalid."
      ]
   },
   "valid":false
}
```
