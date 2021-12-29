.PHONY: build maria

maria:
	podman run -p 127.0.0.1:3306:3306  --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest

mongodb:
	podman run -p 127.0.0.1:27017:27017 --name some-mongo -e MONGO_INITDB_ROOT_PASSWORD=secret -e MONGO_INITDB_ROOT_USERNAME=mongoadmin -d mongo:latest

offline-image:
	podman build -t offline-order:latest -f Dockerfile.offline .
online-image:
	podman build -t online-order:latest -f Dockerfile.online .

offline-container:
	podman run -p:9000:9000 --env-file ./offline.env --link some-mariadb:db \
	--rm --name offlineapp offline-order:latest
online-container:
	podman run -p:9001:9001 --env-file ./online.env --link some-mongo:db \
	--rm --name onlineapp online-order:latest
