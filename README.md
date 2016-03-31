# goland

Trying to support the creation of a catalog with different types of items. These items can be retrieved, posted and updated throughout RESTful APIs.

Pending/In progress


*1)* Change GET /catalog/products response to return an array of items but within a named object.
     e.g. 
            {
                "items":[
                          {
                          "id":"1",
                          "descriptions":"blah",
                          "price":1
                          },
                          {
                          "id":"2",
                          "descriptions":"blahblah",
                          "price":2
                          }
                ]
            }
*2)* Android floating button. Try to use this button to be over the screen an no to scroll up or down.

*3)* Support sending pictures when POSTing items.

*4)* Support returning items withing a date range.