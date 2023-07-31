To Run Docker:

docker-compose up --build (for the first run)

docker-compose up (for future runs without rebuilding)

docker-compose up --build --remove-orphans (to rebuild)

Curl function:

curl --location --request POST 'http://127.0.0.1:8080/hello' --header 'Content-Type: application/json' --data-raw '{"message": "Arbitrary Name or Value"}'    

curl --location --request POST 'http://127.0.0.1:8080/comment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "message": "Ashley",
    "comment": "hello there"
}'

-remove-orphans
For Windows:

Invoke-WebRequest -Uri 'http://127.0.0.1:8080/hello' -Method POST -Headers @{'Content-Type'='application/json'} -Body '{"message": "hohohohoho"}'

# Getting Started with Orbital5
This document guides you through the installation of required dependencies, project setup, and using the basic functionalities of Orbital5.

## Prerequisites
To work with Orbital5, you need to install Go, Etcd, and Docker. Follow the instructions below to install these prerequisites.

### Go Installation
1. Download the Go binary from the official Go [download page](https://golang.org/dl/).
2. Extract the binary and set your $PATH to include the Go bin directory. For example:
shell
tar -xvf go1.16.4.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin
3. Verify the Go installation:
shell
go version
You should see the Go version if it's installed successfully.

### Etcd Installation
1. Download the latest version of etcd from the [official Github repository](https://github.com/etcd-io/etcd/releases).
2. Extract the tarball:
shell
tar -xvf etcd-v3.4.13-linux-amd64.tar.gz
3. Move the executable files to /usr/local/bin:
shell
sudo mv etcd-v3.4.13-linux-amd64/etcd* /usr/local/bin/
4. Verify the etcd installation:
shell
etcd --version
This should return the version of etcd if it's installed correctly.

### Docker Installation
1. Update your existing list of packages:
shell
sudo apt update
2. Install a few prerequisite packages which let apt use packages over HTTPS:
shell
sudo apt install apt-transport-https ca-certificates curl software-properties-common
3. Add the GPG key for the official Docker repository to your system:
shell
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
4. Add the Docker repository to APT sources:
shell
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
5. Update the package database with the Docker packages from the newly added repository:
shell
sudo apt update
6. Finally, install Docker:
shell
sudo apt install docker-ce
7. Verify Docker installation:
shell
sudo systemctl status docker
The output should display Docker status as active (running).

## Setting up Orbital5

1. Clone the project from GitHub:
shell
git clone https://github.com/shamesjen/orbital5.git
2. Navigate to the project directory:
shell
cd orbital5
3. Initialize the project using Docker:
shell
docker-compose up --build
At this point, the Orbital5 service is ready to be used!

## Accessing Orbital5 Endpoints

You can test the different methods of Orbital5 using curl commands as follows:

1. The hello method:
shell
curl --location --request POST 'http://127.0.0.1:8080/hello' \
--header 'Content-Type: application/json' \
--data-raw '{"message": "Arbitrary Name or Value"}'

Expected output: {"message":"Hello!, Arbitrary Name or Value"}

2. The like method:
shell
curl --location --request POST 'http://127.0.0.1:8080/like' \
--header 'Content-Type: application/json' \
--data-raw '{"message": "Arbitrary Name or Value"}'

Expected output: {"message":"Arbitrary Name or Value has successfully liked VideoID: 11234"}% Video ID is the current video the viewer is watching.

3. The unlike method:
shell
curl --location --request POST 'http://127.0.0.1:8080/unlike' \
--header 'Content-Type: application/json' \
--data-raw '{"message": "Arbitrary Name or Value"}'

Expected output: {"message":"Arbitrary Name or Value has successfully unliked VideoID: 11234"}%

4. The comment method:
shell
curl --location --request POST 'http://127.0.0.1:8080/comment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "message": "Arbitrary Name or Value",
    "comment": "Text to be commented on the video"
}'

Expected output: {"message":"Arbitrary Name or Value has commented: \"Text to be commented on the video\" on VideoID: 11234"}%

Please replace "Arbitrary Name or Value" with your own values before running these commands. 

Congratulations! You have now successfully setup and accessed the Orbital5 project.