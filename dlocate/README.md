# dlocate utility

## using command line:

```
# first build:
go build

# to index:
./dlocate -o index -d /home/ahmed/Downloads/csed/networks

# to search:
./dlocate -o search -d /home/ahmed/Downloads/csed/networks -s midterm

# to update:
./dlocate -o update -d /home/ahmed/Downloads/cloud\ computing/

# to clear:
./dlocate -o clear
```

## using docker

first put a test folder in indexFiles and name it `testFolder`

```
# build the project

# we use CompileDaemon to build the project if any file changed
sudo docker-compose up -d

# to index
sudo docker-compose exec dlocate ./dlocate -o index -d ./indexFiles/testFolder/

# to search
sudo docker-compose exec dlocate ./dlocate -o search -d ./indexFiles/testFolder/ -s midterm

# to update
sudo docker-compose exec dlocate ./dlocate -o update -d ./indexFiles/testFolder/

# to clear
sudo docker-compose exec dlocate ./dlocate -o clear
```