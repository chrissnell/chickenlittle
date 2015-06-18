# Team API

## Get list of all teams
**Request**
```
GET /teams
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "teams": [
    {
      "description": "Kings", 
      "members": [
        "arthur"
      ], 
      "name": "kings", 
      "rotation": {
        "description": "none", 
        "frequency": 0, 
        "time": "0001-01-01T00:00:00Z"
      }, 
      "steps": [
        {
          "method": 0, 
          "target": "", 
          "timebefore": 3600
        }
      ]
    }
  ],
  "message": "",
  "error": ""
}
```

## Fetch details for a team
**Request**
```
GET /teams/NAME
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "teams": [
    {
      "description": "Kings", 
      "members": [
        "arthur"
      ], 
      "name": "kings", 
      "rotation": {
        "description": "none", 
        "frequency": 0, 
        "time": "0001-01-01T00:00:00Z"
      }, 
      "steps": [
        {
          "method": 0, 
          "target": "", 
          "timebefore": 3600
        }
      ]
    }
  ],
  "message": "",
  "error": ""
}
```

## Create a new team
**Request**
```
POST /teams

{
  "description": "Kings", 
  "members": [
    "arthur"
  ], 
  "name": "kings", 
  "rotation": {
    "description": "none", 
    "frequenc": 0, 
    "time": 0
  }, 
  "steps": {
    "method": 0, 
    "target": "", 
    "timebefore": 3600
  }
}
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "message": "Team kings created",
  "error": ""
}
```

## Update a team
**Request**
```
PUT /teams/NAME

{
  "description": "Kings and Queens",
  "members": [
    "arthur",
    "lancelot"
  ]
}
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "teams": [
    {
      "name": "kings",
      "description": "Kings and Queens"
    }
  ],
  "message": "Team kings updated",
  "error": ""
}
```

## Delete a team
**Request**
```
DELETE /teams/NAME
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "message": "Team kings deleted",
  "error": ""
}
```

