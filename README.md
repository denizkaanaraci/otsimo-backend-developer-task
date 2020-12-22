# otsimo-backend-developer-task

    Database and server configs can be changed in .env file
    
    build:
        go build -o otsimo-app
    
    run: 
        ./otsimo-app

## docker 

    build (multi-stage):
        docker build -t otsimo-app .
    
    run: 
        docker run -d -p 8080:8080 -e ADDR=0.0.0.0:8080 -e DBURI=mongodb://ip:port --name otsimo-app otsimo-app
        
        interactive: docker run -it --rm -p 8080:8080 -e ADDR=0.0.0.0:8080 -e DBURI=mongodb://ip:port --name otsimo-app otsimo-app

## Functions

### Create Candidate

**URI:** `/candidate`

**Method:** `POST`

**Example:** http://0.0.0.0:8080/candidate

```json
{ 
    "first_name" : "first_name", 
    "last_name" : "last_name", 
    "email" : "test@hotmil.com", 
    "department" : "Design", 
    "university" : "Ankara", 
    "experience" : false
}
```

### Read Candidate

**URI:** `/candidate/{id}`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/candidate/5b758c6151d9590001def630

### Delete Candidate

**URI:** `/candidate/{id}/delete`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/candidate/5b758c6151d9590001def630/delete

### Arrange Meeting

**URI:** `/meeting/arrange`

**Method:** `POST`

**Example:** http://0.0.0.0:8080/meeting/arrange

```json
{
  "_id": "5fe106bdc877325c8450f754",
  "nextMeetingTime": "123456789"
}
```

### Complete Meeting

**URI:** `/meeting/complete/{id}`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/meeting/complete/5b758c6151d9590001def630


### Deny Candidate

**URI:** `/candidate/{id}/deny`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/candidate/5b758c6151d9590001def630/deny

### Accept Candidate

**URI:** `/candidate/{id}/accept`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/candidate/5b758c6151d9590001def630/accept

### Find Assignee ID By Name


**URI:** `/assignee?name={id}`

**Method:** `GET`

**Params:** `id: candidate id`

**Example:** http://0.0.0.0:8080/assignee?name=Zafer

### Find Assignees Candidates

**URI:** `/assignee/{id}/candidates`

**Method:** `GET`

**Params:** `id: assignee id`

**Example:** http://0.0.0.0:8080/assignee/5c191acea7948900011168d4/candidates