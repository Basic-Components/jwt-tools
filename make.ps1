$ASSETS = "bin"
$GOARCHS = "386", "amd64"
$GOOSS = "linux", "darwin", "windows"
$env:GO111MODULE="on"
# Set the GOPROXY environment variable
$env:GOPROXY="https://goproxy.io"


$cmd = "win64"
$name = "jwtcenter"
if ($args.Count -eq 0){
    $cmd = "win64"
}elseif ($args.Count -eq 1){
    $cmd = $args[0]
}elseif ($args.Count -eq 2){
    $cmd = $args[0]
    $name = $args[1]
}else{
    echo "args too much"
    exit
}
 
if (!(Test-Path $ASSETS)) {
    mkdir $ASSETS
    protoc -I=schema --go_out=plugins=grpc:jwtcenter/jwtrpcdeclare --go_out=plugins=grpc:jwtcentersdk/jwtrpcdeclare --go_opt=paths=source_relative jwtrpcdeclare.proto
} 

if ($cmd -eq "all"){
    foreach ($env:GOARCH in $GOARCHS) {
        foreach ($env:GOOS in $GOOSS){
            $target = "$ASSETS/$env:GOOS-$env:GOARCH"
            if (!(Test-Path $target)){
                mkdir $target
            }
            if ($env:GOOS -eq "windows"){
                go build -ldflags "-s -w" -o $target/$name.exe jwtcenter/main.go
            }else {
                go build -ldflags "-s -w" -o $target/$name jwtcenter/main.go
            }
            
        }
    }
}elseif ($cmd -eq "win32") {
    $env:GOARCH="386"
    $env:GOOS="windows"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name.exe jwtcenter/main.go
}elseif ($cmd -eq "win64") {
    $env:GOARCH="amd64"
    $env:GOOS="windows"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name.exe jwtcenter/main.go
}elseif ($cmd -eq "mac") {
    $env:GOARCH="amd64"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name  jwtcenter/main.go
}elseif ($cmd -eq "mac32") {
    $env:GOARCH="386"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name jwtcenter/main.go
}elseif ($cmd -eq "linux32") {
    $env:GOARCH="386"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name jwtcenter/main.go
}elseif ($cmd -eq "linux64") {
    $env:GOARCH="amd64"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name jwtcenter/main.go
}elseif ($cmd -eq "linuxarm") {
    $env:GOARCH="arm"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -ldflags "-s -w" -o $target/$name jwtcenter/main.go
}else{
    echo "unknown cmd $cmd"
}