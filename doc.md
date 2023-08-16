## Technical Documentation

### Completed Functionalities
1. Login
2. Registration
3. Errors


### Bugs
1. Login Bug #1.0
2. Bug #2.0: Both TCP and HTTP server shuts down after receiving a file
3. Bug #3.0: If TCP client closes connection, server keeps reading endless stream of data.


### Design Decisions
#### File Sending
After file data is sent, the client will also send the filemetadata (containing file name, file type, etc) after 100ms. To avoid the server confusing this process with repeatative updates, the client will do the following:  
- After uploading a file change, start a counter that will wait 200ms.
- If the same file is updated again within 200ms, log that and wait until the counter is finished to upload the file