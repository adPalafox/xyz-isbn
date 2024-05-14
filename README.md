# xyz-isbn

### Running the program
Please ensure that you have Golang in your local machine

The web app XYZ Books is needed to be running before executing this program, please make sure that the server is running properly from the other repository.

To run the program:
```
make run-service
```
Alternatively you can also do it directly with:
```
go run main.go
```

### Process
The service will fetch the list of books enpoint, and automatically check for missing ISBN information and perform the necessary steps:
 1. Will make concurrent fetch to GET API of XYZ Books web app to list the paginated books per batch
 2. Check for missing ISBN  13 or 10 information
 3. Convert the missing ISBN to proper format
 4. Will update the Book record from the web app via the PATCH API route
 5. Update the CSV for the ISBN record
