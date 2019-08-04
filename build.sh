bash make.sh linux64
source version
docker build -t hsz1273327/jwtrpc:$PROJECT_VERSION -t hsz1273327/jwtrpc:latest .
docker push hsz1273327/jwtrpc
git add .
git submit -m "update to v$PROJECT_VERSION"
git push
git tag v$PROJECT_VERSION
git push origin --tags