# project-viewMore
Contains a basic backend project to support movie upload and feedback



INSTRUCTION TO RUN BINARY:

1. create a directory "project" in "src" directory present in GOPATH
2. cd /project
3. git clone https://github.com/malishan/project-viewMore.git
4. set the following environment variables:
        export MongoHost="mongodb://localhost:27017"
        export RedisHost="localhost:6379"
5. start mongo and redis server
    mongod --dbpath="./mongodatabase"
    redis-server /usr/local/etc/redis.conf
6. execute - go run main.go     (OR)     ./project-viewMore
7. call postmant APIs present in json files
8. mongo can be loaded with values present in the available json files





POINTS TO NOTE:

1. Login API not functioning it
2. Movie can only be added by admin
3. Rating and Comment can only be given by loggedIn users
4. Movie Search can be done by anyone
5. Feedback available only for loggedIn users
6. Once a user has rated on a movie, he/she cannot rate again
7. Multiple comments on a movie can be put up by a single loggedIn user




NOTE:

1. All APIs are present in json files (2 files containing different format)
2. Three Mongo Collections are present in json files
3. ACL not handled
4. Facing issue with redis, thus login not implemented completely