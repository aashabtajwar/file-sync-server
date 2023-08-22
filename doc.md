# Technical Documentation

## Completed Functionalities

1. Login
2. Registration
3. Errors

## To Modify

## Bugs

1. Login Bug #1.0
2. Bug #2.0: Both TCP and HTTP server shuts down after receiving a file
3. Bug #3.0: If TCP client closes connection, server keeps reading endless stream of data.
4. Bug #4.0: bytes.Buffer is printing a string?? -_- seriously?  
Update on Bug #4.0:  Fixed (it was not a bug, check. Because a value of type *bytes.Buffer has a String() method (the method set of *bytes.Buffer contains the String() method), and a value of type bytes.Buffer does not. And the fmt package checks if the value being printed has a String() string method, and if so, it is called to produce the string representation of the value.)

## Design Decisions

## File Sending

After file data is sent, the client will also send the filemetadata (containing file name, file type, etc) after 100ms. To avoid the server confusing this process with repeatative updates, the client will do the following:  

- After uploading a file change, start a counter that will wait 200ms.
- If the same file is updated again within 200ms, log that and wait until the counter is finished to upload the file

## Issues

1. How to handle file rename?
