# dlocate utility

## using command line:

```
# first build:
go build

# to index:
./dlocate -o index -d /home/ahmed/Downloads/csed/networks

# to search in file names:
./dlocate -o searchNames -d /home/ahmed/Downloads/csed/networks -s midterm

# to search in file content:
./dlocate -o searchContent -d /home/ahmed/Downloads/csed/networks -s midterm


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
run the project
sudo docker-compose up -d

# to index
sudo docker-compose exec dlocate ./dlocate -o index -d ./indexFiles/testFolder/

# to search in file names
sudo docker-compose exec dlocate ./dlocate -o searchNames -d ./indexFiles/testFolder/ -s midterm

# to search in file content
sudo docker-compose exec dlocate ./dlocate -o searchContent -d ./indexFiles/testFolder/ -s midterm

# to update
sudo docker-compose exec dlocate ./dlocate -o update -d ./indexFiles/testFolder/

# to clear
sudo docker-compose exec dlocate ./dlocate -o clear
```