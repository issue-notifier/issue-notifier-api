# Issue Notifier API

The API service for the Issue Notifier website.

### Feature Sets (yet to be implemented)
- [ ] An API to fetch latest issues per repository on demand which will be called from the UI
- [ ] Set of APIs to manage user preferences for frequency of emails

Feel free to raise PRs for the above mentioned features or you can also raise issues if you think you have a new feature request.

### To run the service locally

#### Without Minikube
1. You need to have Go & PostgreSQL installed
2. Setup env vars
3. Setup the database. Run the commands in `database.txt` file from psql command line
3. Run `$ go run main.go` 

Swagger Page - http://localhost:8001/api/v1/swagger/index.html

#### With Minikube
1. You need to have docker, minikube, [kubeseal](https://tomcode.com/blog/how-to-manage-kubernetes-secrets-securely-in-git), Go & PostgreSQL installed
2. `$ minikube start`
3. [Register a new OAuth app](https://github.com/settings/applications/new) in GitHub. Add Homepage URL and Authorization callback URL as `http://<minikube ip>`. You can get your minikube ip by running the command 
`$ minikube ip`
Use the Authorization callback URL, Client ID and Client Secret below for creating secrets and configmap

Run the following commands from root of the project

4. `$ eval $(minikube docker-env)`
5. `$ docker build -t hemakshis/issue-notifier-api .` (also make sure to build hemakshi/issue-notifier image after this, details here - [issue-notifier/issue-notifier](https://github.com/issue-notifier/issue-notifier))
6. `$ cd deploy/development`
7. `$ kubectl apply -f namespace.yaml`
8. `$ kubectl -n issue-notifier create secret generic secrets --from-literal=dbPass=<pass> --from-literal=githubClientID=<client_id> --from-literal=githubClientSecret=<client_secret> --from-literal=sessionAuthKey=<auth_key> --dry-run -o yaml > secrets.yaml`
9. `$ kubeseal --controller-namespace=<kubeseal_namespace> --format=yaml < secrets.yaml > sealed-secrets.yaml`
10. Update `comfig-map.yaml` with your own Authorization callback URL and Client ID
11. Finally run `$ kubectl apply -f sealed-secrets.yaml -f config-map.yaml -f postgres.yaml -f issue-notifier.yaml -f issue-notifier-api.yaml -f ingress-service.yaml`
12. Create database and tables 
`$ kubectl -n issue-notifier exec -it <postgres_pod> -- /bin/bash`
`$ psql -h localhost -p 5432 -U postgres`
Run the commands in database.txt file
13. In your browser go to [http://<minikube_ip>/api/v1/health]() and you should get a response of `{"status":"UP"}`

### Contribution
1. Raise a bug or a feature request
2. Keep checking the Issues tab.
3. Find & solve `TODO`s in the source code and raise a PR
4. You can write unit tests!

#### Contact
Reach out to [Hemakshi Sachdev](https://github.com/hemakshis) for any queries.