# Notification API

### Notify a person

**Request**
```
POST /people/USERNAME/notify

    {
        "content": "Dinnertime, chickies, lets all eat.  Wash your wings and take a seat.",
    }
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "username": "USERNAME",
  "uuid": "d6b65a80-5a58-4334-8f25-c35619998ba5",
  "content": "Dinnertime, chickies, lets all eat.  Wash your wings and take a seat.",
  "message": "Notification initiated",
  "error": ""
}
```

### Stop an in-progress notification

**Request**
```
DELETE /notifications/UUID
```

**Example Response**
```
HTTP/1.1 200 OK
```
```json
{
  "username": "",
  "uuid": "81ce4c82-6e78-4491-9fbe-537bdce4459a",
  "content": "",
  "message": "Attempting to terminate notification",
  "error": ""
}
```
