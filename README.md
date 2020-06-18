# DeepLocate
Advanced File System Search Engine that have smart ways for search specified for the file system and use deep learning search in the file contents


# dlocate utility

## using command line:

```
# first build:
go build

# to index:
./dlocate -o index -d /home/ahmed/Downloads/csed/networks

optional: -deepScan to use ML for extacting files content

# to open the web
./dlocate -o web

# to search in file names:
./dlocate -o search -d /home/ahmed/Downloads/csed/networks -s midterm

# to search in file content:
./dlocate -o search -deepScan -d /home/ahmed/Downloads/csed/networks -s midterm

# to search in meatadata:
./dlocate -o metaSearch -d /home/ahmed/Downloadscsed/networks/ -s clock --deepScan

# to update:
./dlocate -o update -d /home/ahmed/Downloads/csed/networks

optional: -deepScan to use ML for extacting files content

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

# to open the web
sudo docker-compose exec dlocate ./dlocate -o web


# and so on like the previous block ...
```