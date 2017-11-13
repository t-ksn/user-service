 branch | status| coverage |
 -------|-------|----------|
 master| [![Build Status](https://travis-ci.org/t-ksn/user-service.svg?branch=master)](https://travis-ci.org/t-ksn/user-service)| [![codecov](https://codecov.io/gh/t-ksn/user-service/branch/master/graph/badge.svg)](https://codecov.io/gh/t-ksn/user-service)
# User service
The User service provide API for manage user account and join/leave group and company.

# Table of content
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
port | PORT | Listening port (`by default 8080`)

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
102 - Minimum password length 4
103 - User name is empty
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
