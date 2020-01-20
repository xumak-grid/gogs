.PHONY: deploy deploy deploy-service remove secret

deploy-test:
	kubectl create -f k8s/test.yml

deploy-service:
	kubectl --namespace=ehernandez create -f k8s/gogs-service.yml

deploy-pod:
	kubectl --namespace=ehernandez create -f k8s/gogs.yml
	
remove:
	kubectl --namespace=ehernandez delete -f k8s/gogs.yml || true
	kubectl --namespace=ehernandez delete -f k8s/gogs-service.yml || true
	kubectl --namespace=ehernandez delete -f k8s/test.yml || true

secret:
	kubectl delete secret regsecret
	kubectl --namespace=ehernandez --kubeconfig=../kubeconfig create secret docker-registry regsecret --docker-server=https://598240822331.dkr.ecr.us-east-1.amazonaws.com --docker-username=AWS --docker-password="$AWS_ECR_LOGIN" --docker-email=none

new-tag:
	git tag -a "$(VERSION)" -m "version $(VERSION)"
	git push --tags
