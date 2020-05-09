.PHONY: create-secret create-ns deploy download-secrets create-configmap delete-ns
.SILENT: create-secret

NS ?= authentic

certificate.yaml:
	echo "Download the certificate.yaml before running the code"

firebase-web-apikey.yaml:
	echo "Download the firebase-web-apikey.yaml before running the code"

create-ns:
	kubectl create ns ${NS}

create-configmap:
	kubectl create --save-config configmap nginx-conf.d --from-file=nginx-conf.d --dry-run -o yaml > config.yaml
	kubectl apply -f config.yaml -n ${NS}

deploy-secrets: certificate.yaml firebase-web-apikey.yaml
	kubectl apply -f certificate.yaml -n authentic
	kubectl apply -f firebase-web-apikey.yaml -n authentic

deploy: create-configmap
	kubectl apply -f Resources.yaml -n authentic

create-secret:
	kubectl create secret generic firebase-web-apikey \
		--from-literal=FIREBASE_WEB_APIKEY=${FIREBASE_WEB_APIKEY} \
		--dry-run -o yaml \
		-n authentic > firebase-web-apikey.yaml

download-secrets:
	kubectl get secret certificate -n ${NS} -o yaml > certificate.yaml
	kubectl get secret firebase-web-apikey -n ${NS} -o yaml > firebase-web-apikey.yaml

delete-ns:
	kubectl delete ns ${NS}