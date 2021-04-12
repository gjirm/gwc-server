#!/usr/bin/pwsh

$projectName = "Go Wireguard Control server"
$dt = Get-Date -Format "yyy-MM-dd_HHMMss"
$headhash = git rev-parse --short HEAD
$tag = "$($dt)-$($headhash)"
$imageName = "jirm/gwc-server:$($tag)"
$imageLatest = "jirm/gwc-server:latest"

$minisignKey = "W:\keys\jirm-minisign-2020.key"

Write-Host "--> $projectName <--" -ForegroundColor Green
if ($Args[0] -eq "build") {
    if ($Args[1] -eq "docker") {

        Write-Host "--> Building $($imageName)" -ForegroundColor Green
        docker build --tag $imageName --tag $imageLatest .
        If ($lastExitCode -eq "0") {
            Write-Host "--> $($imageName) successfully build!" -ForegroundColor Green
        } else {
            Write-Host "--X $($imageName) build failed!" -ForegroundColor Red
        }
    } elseif ($args[1] -eq "cli") {

        Write-Host "--> Building WebWormhole CLI version $tag" -ForegroundColor Green
        go mod download
        go build -o ww.exe .\cmd\ww
        Write-Host "--> Building CLI" -ForegroundColor Green
        minisign -Sm ww.exe -s $minisignKey -c "WebWormhole CLI version $tag - signed $(Split-Path -Leaf $minisignKey)" -t "WebWormhole CLI version $tag - signed $(Split-Path -Leaf $minisignKey)"
    } else {

        Write-Host "None!" -ForegroundColor Yellow
    }
} elseif ($args[0] -eq "run") {

    if ($args[1] -eq "docker") {

        Write-Host "--> Run Docker container"  -ForegroundColor Green
        docker run --rm -v ${PWD}/config.yml:/gwc/config.yml -v ${PWD}/log:/gwc/log/ -v ${PWD}/db:/gwc/db/ -p 8080:8080 jirm/gwc-server
    }
} else {
    Write-Host "--! None!" -ForegroundColor Yellow
}
