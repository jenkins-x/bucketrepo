CHART_REPO := http://jenkins-x-chartmuseum:8080
NAME := jx-app-bucketrepo

init: 
	helm init --client-only

setup: init
	helm repo add jenkins-x https://storage.googleapis.com/chartmuseum.jenkins-x.io

build: clean setup
	jx step helm build
	helm lint
	helm template .

install: clean setup build
	helm install . --name ${NAME}

upgrade: clean setup build
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
	curl --fail -u $(CHARTMUSEUM_USER):$(CHARTMUSEUM_PASS) --data-binary "@$(NAME)-$(VERSION).tgz" $(CHART_REPO)/api/charts
