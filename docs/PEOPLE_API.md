# People API

## Get list of all people
**Request**
```
GET /people
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "people": [
    {
      "username": "arthur",
      "fullname": "King Arthur"
    },
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot"
    }
  ],
  "message": "",
  "error": ""
}
```

## Fetch details for a person
**Request**
```
GET /people/USERNAME
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "people": [
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot"
    }
  ],
  "message": "",
  "error": ""
}
```

## Create a new person
**Request**
```
POST /people

{
  "username": "lancelot",
  "fullname": "Sir Lancelot"
}
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "message": "User lancelot created",
  "error": ""
}
```

## Update a person
**Request**
```
PUT /people/USERNAME

{
  "fullname": "Sir Lancelot the Brave"
}
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "people": [
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot the Brave"
    }
  ],
  "message": "User lancelot updated",
  "error": ""
}
```

## Delete a person
**Request**
```
DELETE /people/USERNAME
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "message": "User lancelot deleted",
  "error": ""
}
```
