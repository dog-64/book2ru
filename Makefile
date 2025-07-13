BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

NAME   := dockerhub.duletsky.ru/novik_app
TAG    := $$(git log -1 --pretty=%!H(MISSING))
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

ARGS = $(filter-out $@,$(MAKECMDGOALS))
%:
	@:
.PHONY: spec test

build:
	echo $(BRANCH)
	time go build

b:
	make build
	docker-compose images

r:
	./drun.sh

test:
	bundle
	make rspec
	make rubocop
	# make fasterer
	bundle exec bundle-audit --update
	make brakeman
	trivy fs . --scanners vuln --skip-dirs ./vendor
	go test

rspec:
	#bundle exec rspec
	#RAILS_ENV=test bundle exec rspec  --format d
	RAILS_ENV=test bundle exec rspec

rubocop:
	# bundle exec  rubocop -A -F ./app ./lib ./spec
	bundle exec  rubocop -A

brakeman:
	bundle exec brakeman -q -6 --no-summary --force

