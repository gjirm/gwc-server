#!/usr/bin/pwsh

$projectName = "Go Wireguard Control server"
$dt = Get-Date -Format "yyy-MM-dd_HHMMss"
$headhash = git rev-parse --short HEAD
$tag = git describe --tags --abbrev=0
$imageName  = "jirm/gwc-server:$($headhash)"
$imageNameTag = "jirm/gwc-server:$($tag)"
$imageLatest = "jirm/gwc-server:latest"

$minisignKey = "W:\keys\jirm-minisign-2020.key"

Write-Host "--> $projectName <--" -ForegroundColor Green

if ($Args[0] -eq "build-docker-no-cache-tag") {

    Write-Host "--> Building $($imageName)" -ForegroundColor Green
    docker build --no-cache --tag $imageNameTag --tag $imageName --tag $imageLatest .
    If ($lastExitCode -eq "0") {
        Write-Host "--> $($imageName) successfully build!" -ForegroundColor Green
    } else {
        Write-Host "--X $($imageName) build failed!" -ForegroundColor Red
    }
}

if ($args[0] -eq "run-docker") {

    Write-Host "--> Run Docker container"  -ForegroundColor Green
    docker run --rm -v $PSScriptRoot\config.yml:/gwc/config.yml --name gwc -p 8080:8080 jirm/gwc-server
}

if ($args[0] -eq "stop-docker") {

    Write-Host "--> Stop Docker container"  -ForegroundColor Green
    docker stop gwc
}


Write-Host "--! None!" -ForegroundColor Yellow


# Write-Host "--> Building WebWormhole CLI version $tag" -ForegroundColor Green
# go mod download
# go build -o ww.exe .\cmd\ww
# Write-Host "--> Building CLI" -ForegroundColor Green
# minisign -Sm ww.exe -s $minisignKey -c "WebWormhole CLI version $tag - signed $(Split-Path -Leaf $minisignKey)" -t "WebWormhole CLI version $tag - signed $(Split-Path -Leaf $minisignKey)"
