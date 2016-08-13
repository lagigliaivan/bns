# goland

Thi project is intended to be the backend for an android app. It allows
recording purchases by month. Later this can be used to create a user profile and then
recommend where to buy a product a its best price. 

As a whole, these two applications allow keeping records of different purchases and some statistics about how much money was spent by month.

# Compiling bns
/opt/go/bin/go build -o /tmp/pos /home/ivan/dev/bns/src/services/main.go

### Creating docker image
cd /home/ivan/dev/bns/src/services
sudo docker build -t bns/white:0.0.4 .
sudo docker run -d -p 8080:8080 -e DB_TYPE=LOCALDB -e DB_URL=http://192.168.0.5:8000 -e USER_VALIDATION=NON_G_USER bns/white:0.0.4


**Using remote dynamo (if you are running dynamo in your machine, then AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are not needed)**
sudo docker run -d -p 8080:8080 -e USER_VALIDATION=NON_G_USER -e AWS_ACCESS_KEY_ID=<Key provided to you by amazon> -e AWS_SECRET_ACCESS_KEY=<secret provided to you by amazon> bns/white:0.0.4

Note: Once the docker example above is run, you can use 
curl -H "Authorization:d563af2d08b4f672a11b3ed9065b7890a6412cab" http://localhost:8080/catalog/purchases
OR
curl -H "Authorization:d563af2d08b4f672a11b3ed9065b7890a6412cab" http://localhost:8080/catalog/purchases?groupBy=month

When running the docker image some environment variables can be used to use it locally.

DB_TYPE: [LOCALDB, MEMDB]
DB_URL: If LOCALDB is used, then you can use this variable to set the dynamo db url. localhost:8000 by default, if this variable is not used.
ANDROID_APP_ID : android app id to validate user token
USER_VALIDATION : [NON_G_USER] 



### Docker container for Local AWS DynamoDB

AWS DynamoDB Local will let you test against DynamoDB without needing
a full network. For details see https://aws.amazon.com/blogs/aws/dynamodb-local-for-desktop-development/

To use link to your application:
sudo docker run -d --name dynamodb deangiberson/aws-dynamodb-local   *NO EXPORTING PORT IS NEEDED SINCE YOU CAN USE THE IP ADDRESS OF THE DOCKER CONTAINER"
sudo docker run -d -P --name web --link dynamodb:dynamodb training/webapp python app.py


### Creating a table in a local dynamoDB
aws dynamodb create-table --table-name Purchases --attribute-definitions AttributeName=id,AttributeType=S AttributeName=dt,AttributeType=N --key-schema AttributeName=id,KeyType=HASH AttributeName=dt,KeyType=range --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url  http://localhost:8000
aws dynamodb create-table --table-name Items --attribute-definitions AttributeName=user_purchase,AttributeType=S AttributeName=item,AttributeType=S --key-schema AttributeName=user_purchase,KeyType=HASH AttributeName=item,KeyType=range 
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://172.17.0.2:8000

### Listing tables in a local dynamoDB

aws dynamodb list-tables --endpoint-url http://<dockerip>:8000

### Deleting a table
aws dynamodb delete-table --table-name Purchases --endpoint-url http://172.17.0.2:8000


