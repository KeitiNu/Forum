FROM golang:1.17-buster

#Selects everything in current directory and copies
ADD ./ /forumcontainer

# Move to working directory /forumcontainer/cmd/web
WORKDIR /forumcontainer/cmd/web

# Build the application
RUN go build --tags "sqlite_userauth" -o forumapp

# Move to working directory /forumcontainer
WORKDIR /forumcontainer

# Export necessary port
EXPOSE 8090

# Command to run when starting the container
CMD ["./cmd/web/forumapp"]
