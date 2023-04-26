main:
	docker-compose up --build --force-recreate --no-deps

deploy:
	export AWS_PROFILE=private
	aws ecr get-login-password --region eu-north-1 | docker login --username AWS --password-stdin 886937713965.dkr.ecr.eu-north-1.amazonaws.com
	docker build . -t nbesschetnov/challenge
	docker tag nbesschetnov/challenge 886937713965.dkr.ecr.eu-north-1.amazonaws.com/challenge:latest
	docker push 886937713965.dkr.ecr.eu-north-1.amazonaws.com/challenge:latest
	kubectl rollout restart deployment hello-eks-a
