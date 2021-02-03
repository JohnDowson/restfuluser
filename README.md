# restfuluser
Simple RESTful CRUD app using JSON file for storage.

## Usage
* `curl -X GET /user/` -> Returns a list of all users
* `curl -X GET /user/:id` -> Returns a user with matching ID
* `curl -H "Content-Type: application/json" -X POST -d "{ "name": "Username" }"  /user` -> Creates a new user
* `curl -H "Content-Type: application/json" -X PUT -d "{ "name": "Username" }"  /user/:id` -> Updates a user with matching ID with given data
* `curl -X DELETE /user/:id` -> Deletes a user with matching ID
