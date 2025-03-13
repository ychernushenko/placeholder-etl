## Install
brew install docker-compose
brew install go

## Start
docker-compose down --volumes --rmi all
docker-compose build --no-cache
docker-compose up

## Debug
docker exec -it postgres bash
psql -U postgres -d placeholder_db
select * from posts;