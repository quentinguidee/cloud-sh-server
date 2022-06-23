# Routes documentation

## Login

Returns URL to login with GitHub.

`GET /auth/github/login`

- Parameters: `None`
- Data Parameters: `None`
- Success Response: `200`
- Error Response: `None`
- Returned Content:

```json
{
  "url": "https://github.com/login/oauth/authorize?access_type=offline&client_id=client_id&redirect_uri=redirect_uri&response_type=code&scope=all&state=state"
}
```

## Authorize

Validates the code, add the user to the database, and returns the app token and the user.

`GET /auth/github/authorize`

- Parameters: `None`
- Data Parameters:

```json
{
  "code": "code",
  "state": "state"
}
```

- Success Response: `200`
- Error Responses:

  - `400`:
    - Missing parameter(s)
    - The code is invalid
    - The state parameter does not match the one in the URL
  - `500`:
    - Server failed to generate the app token
    - Server failed to retrieve the user from Github
    - Server failed decode Github user account
  - `Any other error`
    - Github related error, see [Github User Api Documentation](https://docs.github.com/en/rest/users/users#get-the-authenticated-user)

- Returned Content:

```json
{
  "token": "token"
}
```

## User

Returns user information.

`GET /auth/user/username`

- Parameters: `username` (string) - username of user to get
- Data Parameters: `None`
- Success Response: `200`
- Error Responses:
  - `404`:
    - User not found
  - `500`:
    - Server failed to retrieve user
- Returned Content:

```json
{
  "name": "John Doe",
  "id": "1",
  "username": "johndoe"
}
```
