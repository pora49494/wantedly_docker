docker build -t go_web ./go_web
docker build -t go_db ./go_db
docker-compose up -d 
docker exec db /etc/init.d/postgresql start
docker exec db /bin/bash /usr/local/bin/db_setup.sh