# User service
The User service provide API for manage user account and join/leave group and company.

<!-- MDTOC maxdepth:6 firsth1:2 numbering:0 flatten:0 bullets:1 updateOnSave:1 -->

- [Dependencies](#dependencies)   
- [Settings](#settings)   
- [Endpoints](#endpoints)   
   - [Register user account](#register-user-account)   
      - [Request](#request)   
      - [Response](#response)   
      - [Errors](#errors)   
   - [Authorization](#authorization)   
      - [Request](#request)   
      - [Response](#response)   
      - [Errors](#errors)   

<!-- /MDTOC -->

## Dependencies
* MongoDB 3
* Go v1.8

## Settings
CMD args (prefix: --)| Env name | Description
---|---|---
db | DB | Connection string to MongoDB (`required`)|
key | KEY | Secret key for sign JWT (`required`) |
port | PORT | Listening port (`by default 80`)

## Endpoints
### Register user account
Create user account

#### Request
**POST /account**
```json
{
    "name": "user_name",
    "password":"user_password"
}
```
#### Response
Success result return 201

#### Errors
```
100 - User name already exist
```

### Authorization
Authorizate user in the system

#### Request
**POST /login**
```json
{
    "name": "user_name",
    "password":"user_password"
}
```
#### Response
```json
{
  "token": "access_token",
  "refresh": "refresh_token",
  "type": "token_type"
}
```
#### Errors
```
101 - User name or password incorrect
```
