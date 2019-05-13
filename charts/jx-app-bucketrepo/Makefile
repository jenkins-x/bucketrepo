CHART_REPO := https://chartmuseum:8080
NAME := jx-app-bucketrepo
CHARTMUSEUM_CREDS_USR := $(shell cat /builder/home/basic-auth-user 2> /dev/null)
CHARTMUSEUM_CREDS_PSW := $(shell cat /builder/home/basic-auth-pass 2> /dev/null)

setup:
	minikube addons enable ingress
	brew install kubernetes-helm
	helm init

build: clean
	jx step helm build
	helm lint
	helm template .

install: clean build
	helm install . --name ${NAME}

upgrade: clean build
	helm upgrade ${NAME} .

delete:
	helm delete --purge ${NAME}

clean:
	rm -rf requirements.lock
	rm -rf charts
	rm -rf ${NAME}*.tgz

release: clean
	sed -i -e "s/version:.*/version: $(VERSION)/" Chart.yaml
	echo "  version: $(VERSION)" >> requirements.yaml
	jx step helm build
	helm package .
	curl --fail -u $(CHARTMUSEUM_CREDS_USR):$(CHARTMUSEUM_CREDS_PSW) --data-binary "@$(NAME)-$(VERSION).tgz" $(CHART_REPO)/api/charts
