./make.ps1 linux64
$env:PROJECT_VERSION="2.0.0"
docker build -t hsz1273327/jwtcenter:$env:PROJECT_VERSION -t hsz1273327/jwtcenter:latest .
docker push hsz1273327/jwtcenter
# git add .
# git commit -m "update to v$env:PROJECT_VERSION"
# git push
# git tag v$env:PROJECT_VERSION
# git push origin --tags