SUDO_DOCKER ?=

all: help

help:
	@echo "make test - run tests (api and cli)"
	@echo "make test-api - run api tests"
	@echo "make test-cli - run cli tests"
	@echo "make clean - stop and remove test containers"
	@echo "make pull - pull Docker images on registry"
	@echo "make init - lauch ambari infra for test purpose"
	@echo "make build - compile program"

init:
	${SUDO_DOCKER} docker-compose up -d ambari-server
	${SUDO_DOCKER} docker-compose up -d ambari-agent
	${SUDO_DOCKER} docker-compose up -d ambari-agent2
	${SUDO_DOCKER} docker-compose up -d ambari-agent3
	until $$(docker-compose run --rm curl --output /dev/null --silent --head --fail http://ambari-server:8080); do sleep 5; done

test-api: clean init
	${SUDO_DOCKER} docker-compose run --rm test

test-cli: clean init
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-or-update-repository --repository-file /workspace/fixtures/repository.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-or-update-repository --use-spacewalk --repository-file /workspace/fixtures/repository.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-or-update-repository --repository-file /workspace/fixtures/repository.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-cluster-if-not-exist --cluster-name test --blueprint-file /workspace/fixtures/blueprint.json --hosts-template-file /workspace/fixtures/cluster-template.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-cluster-if-not-exist --cluster-name test --blueprint-file /workspace/fixtures/blueprint.json --hosts-template-file /workspace/fixtures/cluster-template.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin create-or-update-privileges --privileges-file /workspace/fixtures/privileges.json
	${SUDO_DOCKER} docker-compose run --rm cli --ambari-url http://ambari-server:8080/api/v1 --ambari-login admin --ambari-password admin --debug configure-kerberos --cluster-name "test" --kdc-type "mit-kdc" --kdc-hosts "kdc.test.local" --realm "TEST.LOCAL" --admin-server-host "kdc.test.local" --principal-name "admin/admin@TEST.LOCAL" --principal-password "adminadmin" --domains "test.local,.test.local"

test: test-api test-cli

build:
	${SUDO_DOCKER} docker-compose run --rm build
	
pull:
	${SUDO_DOCKER} docker-compose pull --ignore-pull-failures

clean:
	${SUDO_DOCKER} docker-compose logs ambari-server || exit 0
	${SUDO_DOCKER} docker-compose logs ambari-agent || exit 0
	${SUDO_DOCKER} docker-compose down -v

release:
	mkdir -p release
	@echo "Do nothink"
