Steps to Run this agent.

Change the server IP in /transport/http.go line no 14

Command to build the EXE : go build /cmd/main.go
command to simply run the agent: go run /cmd/main.go

Running the server API 

create a virtual environment to run the API

command: source venv venv/bin/activate

running the API: uvicorn server:app --host 0.0.0.0 --port 8000

After that  wait for 10 secound.