CHART_REPO := gs://jenkinsxio/charts
NAME := bucketrepo

build: clean
	rm -rf Chart.lock
	helm dependency build
	helm lint

install: clean build
	helm install . --name ${NAME}

upgrade: clean build
	helm upgrade ${NAME} .

delete:
	helm delete --purge ${NAME}

clean:
	rm -rf charts
	rm -rf ${NAME}*.tgz

release: clean
	helm dependency build
	helm lint
	helm package .
	helm repo add jx $(CHART_REPO)
	helm gcs push ${NAME}*.tgz jx --public
	rm -rf ${NAME}*.tgz%