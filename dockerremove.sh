#Stops the container called forumcontainer.
docker stop forumcontainer 

#Deletes the container called forumcontainer.
docker rm forumcontainer 

#Deletes the image called forumimage.
docker image rm forumimage

#Deletes the BaseImage.
docker image rm golang:1.17-buster 

#Shows you that all images are gone
docker images

#Shows you that all containers are gone
docker ps -a