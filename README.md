# chickenlittle
A RESTful service to get ahold of people, quickly.

# Notice
This is not ready for public consumption.  Sorry.

# API Usage Examples

## Get list of people
**Request**
```
GET /people
```

**Example Response**
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
```json
{
  "message": "User lancelot created",
  "error": ""
}
```
