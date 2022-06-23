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

Does nothing at the moment.

`GET /auth/github/authorize`

- Parameters: `None`
- Data Parameters: `None`
- Success Response: `200`
- Error Response: `404` (current state of route)
- Returned Content: `None`

## User

Returns user information.

`GET /auth/user/username`

- Parameters: `username` (string) - username of user to get
- Data Parameters: `None`
- Success Response: `200`
- Error Response: `500` (user not found)
- Returned Content:

```json
{
  "name": "John Doe",
  "id": "1",
  "username": "johndoe"
}
```
