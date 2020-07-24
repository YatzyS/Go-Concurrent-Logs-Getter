# Demo Microservice

- This script is made with following tech-stack:
    - Go

- There are 3 files which functions as follows:
    - logcreator.go - To create log files of required size in required quantity for logs_getter testing purpose. 
    - logs_getter.go - To search for logs in a specific date range in the log files.
    - Logs Getter.docx - Describes the thought process which I had while creating the script 
- This is a simple script which can extract logs in a specific date range from a set of large files in very short time.
- Main concept used here is goroutines.
- Following is the procedure to run the code:
    - Clone the project
    - Goto the project directory
    - If you already have logs file in the format :-
        `2006-01-02T15:04:05.00Z,Log A, LogB, Logc\n`
    Then you can use that file else you can create files using the logcreator.go file
    - Then you can run the logs_getter.go with the command:-
        `go run logs_getter.go -f "From Time" -t "To Time" -i "Log file directory location"`
    - Then it will start printing the logs from the file and at the end it will also print how much time it took to print the logs


