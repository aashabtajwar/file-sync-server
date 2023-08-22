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

### How file versioning naming was accomplished

When a file update comes, we have to save the updated file (or the changes in a file) as a different file (or bytes in minio). So how do we find out which version is it currently at?  
Solution: (This approach is when file is stored in the local FS, NOT in Min.io)  
When the file is received for the first time, create two files; first one is named as workspace_userid_filename_versionCurrent.txt (versionCurrent will be named as this to indicate this file contains version current version number). This one will only contain the current version of the file. And the second file will be workspace_userid_filename_v?.mimetype. This is where the actual data will recide in, and the '?' indicates the current version. When updated file comes, first check if the file exists by seaching the first file name. If so, read the version, increment it by one, and finally write the new file with that new version.

## File Sending

After file data is sent, the client will also send the filemetadata (containing file name, file type, etc) after 100ms. To avoid the server confusing this process with repeatative updates, the client will do the following:  

- After uploading a file change, start a counter that will wait 200ms.
- If the same file is updated again within 200ms, log that and wait until the counter is finished to upload the file

## Issues

1. How to handle file rename?
<!-- Ans: When a file rename action occurs, the watcher will detect it -->
