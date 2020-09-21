service:=thunes
mysql:=thunes/mysql
binPath:=$(shell pwd)/bin/
lastCommitMsg:=$(shell git log HEAD -n 1 --pretty=%B)
lastCommitId:=$(shell git rev-parse HEAD | cut -b 1-8)
lastCommitDay:=$(shell git log HEAD -n 1 --pretty=format:%ad --date=format:'%Y%m%d')

branch:=$(shell git rev-parse --abbrev-ref HEAD)
branch:=$(shell echo $(branch) | sed 's/\//-/g')

all : thunes mysql

thunes :
	cp ./thunes.toml $(binPath)
	docker build -q --tag $(service):$(branch)_$(lastCommitId)_$(lastCommitDay) -f ./thunes.dockerfile ./

mysql: 
	docker build -q --tag $(mysql):$(branch)_$(lastCommitId)_$(lastCommitDay) -f ./mysql.dockerfile ./

run : 
	docker run -d --rm\
		-v $(shell pwd)/mysql:/var/lib/mysql \
		-p 127.0.0.1:3306:3306 \
		-p 127.0.0.1:33060:33060 \
		--name thunes-mysql \
		$(mysql):$(branch)_$(lastCommitId)_$(lastCommitDay) \
		mysqld -u root

	docker run -d --rm\
		-p 127.0.0.1:8099:8099 \
		--name thunes \
		--link thunes-mysql \
		$(service):$(branch)_$(lastCommitId)_$(lastCommitDay)

.PHONY: thunes clean run all mysql
clean: 
	rm -rf $(binPath)
	docker image rm $(service):$(branch)_$(lastCommitId)_$(lastCommitDay)
	docker image rm $(service):$(branch)_$(lastCommitId)_$(lastCommitDay) 
