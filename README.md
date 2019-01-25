# Feed service

## Usage
### Prerequirements:
* Docker and docker-compose are installed
* Internet access
* git installed
### Launch
In order to start app with all needed dependencies next steps should be performed:
* clone repo
* run command `docker-compose up`

After startup application should be accessible via http://{docker-machine-host}:8000/

## Design
![Component diagram](/component-diagram.png)

## Service configuration
All configuration parameters should be passed via environment variables. 

| Variable | Description |
| :---- | :---- |
| DB_HOST  | database host |
| DB_USER | database user |
| DB_USER_PASSWORD | database user passwrod |
| DB_NAME | database name |
| KAFKA_HOSTS | kafka hosts |
| FEEDS_TOPIC_NAME | feeds topic name |
| SERVER_PORT | service port - defaut: 8000 |

## API
| API | User story |
| :---- | :---- |
| `POST /{userName}/feed`         |1. As a user, I need an API that allows me to build an activity feed.|
| `GET  /{userName}/feed`         |2. As a user, I need an API to read my own activity feed that contains my activities.|
| `POST /{userName}/action`       |3. As a user, I need an API to follow a friend's activity feed.|
| `GET  /{userName}/feed/friends` |4. As a user, I need an API to retrieve a feed of all activities of friends that I follow.|
| `GET  /{userName}/action` |5. As a user, I need an API to unfollow a friend's feed.|
| `GET  /{userName}/feed?includeRelated=true` |6. For all activity objects in any read feed API, create an additional related field that contains common friend's action..|

### Pagination
You can add pagination parameters to API requests which retrive feed. 

Example: `{url}?page=2&limit=5`. 

Supported API:
* `GET  /{userName}/feed`
* `GET  /{userName}/feed/friends`

Defaults:

| Parameter | Default value |
| :---- | :---- |
| page  | 1 |
| limit | 10 |
 
## Technologies
| Type | Pick |
| :---- | :---- |
| Language | Golang |
| Message broker | Kafka |
| DB | Postgres |
| Container | Docker |
| Orchestration | Docker-compose |

## TODO actions
* Use Graph DB. In future it's better to move to some graph DB for storing relations instead of Postgres RDB.
* Add authentification. Either implement stuff on feed-service instance, either move authentification logic to gateway and resolve user data on feed-service via request headers.
* Extend HETAEOAS. It's better to have backward links, total count, sort etc... parameters supported.
* Unit-test coverage. Although all functionality is being covered by functional tests, all parts of code should also be covered with unit tests.
* Rich validation. Cover all negative cases with validation and return proper response codes.
* Functional testing refactoring. Substitute `sleep()` with retry for asyncroonous calls.






