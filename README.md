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
`localhost:8080/register`
Parameters:
```json
access_code: {access_code}
email: {email}
password: {password}
username: {username}
```
Body:
`email={email}&username={username}&access_code={access_code}&password={password}`
#### Response:
```json
{
   "session":"{session}",
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

### Login `POST`
#### Request:
`localhost:8080/login`
Parameters:
```json
key: {email|username|session}
password: {password}
```
Body:
`key={email|username|session}&password={password}`
#### Response:
```json
{
   "session":"{session}",
   "valid":true
}
```
```json
{
   "valid":false
}
```


