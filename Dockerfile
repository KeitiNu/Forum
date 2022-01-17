FROM golang:1.17-buster

#Selects everything in current directory and copies 
ADD ./ /forumcontainer

# Move to working directory /asciicontainer
WORKDIR /forumcontainer

# Build the application
RUN go build -o forumapp

# Export necessary port
EXPOSE 8090

# Command to run when starting the container
CMD ["./forumapp"]
