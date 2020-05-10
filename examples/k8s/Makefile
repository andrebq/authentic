.PHONY: create-firebase-apikey create-ns deploy download-secrets create-configmap delete-ns restart-nginx
.SILENT: create-firebase-apikey

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
	kubectl apply -f certificate.yaml -n ${NS}
	kubectl apply -f firebase-web-apikey.yaml -n ${NS}

deploy: create-configmap
	kubectl apply -f Resources.yaml -n ${NS}
	kubectl apply -f Daemonset.yaml -n ${NS}

create-firebase-apikey:
	kubectl create secret generic firebase-web-apikey \
		--from-literal=FIREBASE_WEB_APIKEY=${FIREBASE_WEB_APIKEY} \
		--dry-run -o yaml \
		--save-config \
		-n ${NS} > firebase-web-apikey.yaml

download-secrets:
	kubectl get secret certificate -n ${NS} -o yaml > certificate.yaml
	kubectl get secret firebase-web-apikey -n ${NS} -o yaml > firebase-web-apikey.yaml

delete-ns:
	kubectl delete ns ${NS}

restart-nginx: create-configmap
	kubectl scale deployment nginx -n ${NS} --replicas=0
	kubectl scale deployment nginx -n ${NS} --replicas=1

restart-authentic: authentic
	kubectl scale deployment nginx -n ${NS} --replicas=0
	kubectl scale deployment nginx -n ${NS} --replicas=1
