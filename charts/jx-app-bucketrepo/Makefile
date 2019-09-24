CHART_REPO := https://chartmuseum:8080
NAME := jx-app-bucketrepo

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
	curl --fail -u $(CHARTMUSEUM_USER):$(CHARTMUSEUM_PASS) --data-binary "@$(NAME)-$(VERSION).tgz" $(CHART_REPO)/api/charts
