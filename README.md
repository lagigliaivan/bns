# goland

Thi project is intended to be the backend for an android app, which allows
recording purchases by month.

As a whole, these two applications allow keeping records of different purchases and some statistics about how much money was spent by month.

The main idea is to provide some suggestions about where to buy at the best price, the most purchased products by the user, according to his/her
profile. 

### Docker container for Local AWS DynamoDB

AWS DynamoDB Local will let you test against DynamoDB without needing
a full network. For details see https://aws.amazon.com/blogs/aws/dynamodb-local-for-desktop-development/

To use link to your application:
sudo docker run -d --name dynamodb deangiberson/aws-dynamodb-local
sudo docker run -d -P --name web --link dynamodb:dynamodb training/webapp python app.py


### Creating a table in a local dynamoDB
aws dynamodb create-table --table-name Purchases --attribute-definitions AttributeName=user,AttributeType=S
AttributeName=purchase,AttributeType=S --key-schema AttributeName=user,KeyType=HASH AttributeName=purchase,KeyType=range --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://172.17.0.2:8000

aws dynamodb create-table --table-name Items --attribute-definitions AttributeName=user_purchase,AttributeType=S AttributeName=item,AttributeType=S --key-schema AttributeName=user_purchase,KeyType=HASH AttributeName=item,KeyType=range 
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://172.17.0.2:8000

### Listing tables in a local dynamoDB

aws dynamodb list-tables --endpoint-url http://<dockerip>:8000
