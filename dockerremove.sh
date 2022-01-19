#Stops the container called forumcontainer.
docker stop forumcontainer 

#Deletes all the containers not running
docker container prune

#Deletes all the images not used by a container
docker image prune -a

#Shows you that all images are gone
docker images

#Shows you that all containers are gone
docker ps -a