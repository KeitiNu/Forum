
#Builds an image from Dockerfile called forumimage.
docker image build -f Dockerfile -t forumimage .

#Displays all the images on your system.
docker images

#Runs an image called forumimage in a container called forumcontainer in 8090port.
docker container run -p 8090:8090 --detach --name forumcontainer forumimage

#Displays all the containers on your system.
docker ps -a



