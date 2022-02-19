
$CWD = $(Get-Location)
$FRONT = "$CWD\frontend"
$DIST = "$CWD\dist"
$VERSION = Get-Content "$CWD\version.txt"

$NPM_EXEC = $null

function Test-Command( [string] $CommandName ) {
    $null -ne (Get-Command $CommandName -ErrorAction SilentlyContinue)
}

if (!$(Test-Command "7z")) {
    Write-Output "please install 7zip first"
    Write-Output ""
    Write-Output "  Suggest to use Scoop(https://scoop.sh/) to install 7zip (scoop install 7zip)"
    Write-Output ""
    exit
}

if (Test-Command "pnpm") {
    $NPM_EXEC = "pnpm"
}
elseif (Test-Command "yarn") {
    $NPM_EXEC = "yarn"
}
elseif (Test-Command "npm") {
    $NPM_EXEC = "npm"
}

if ($null -eq $NPM_EXEC) {
    Write-Output "no package manager(pnpm, yarn, npm) or nodejs is installled"
    Write-Output ""
    Write-Output "  Suggest to use Scoop(https://scoop.sh/) or NVM to install nodejs (scoop install nodejs),"
    Write-Output "  and use pnpm to manage node packages (npm i -g pnpm)."
    Write-Output ""
}

Remove-Item "$DIST" -Recurse -ErrorAction SilentlyContinue

go build -o "$DIST\daily_wallpaper\daily_wallpaper.exe" -ldflags="-H windowsgui"

Set-Location $FRONT

& $NPM_EXEC install

& $NPM_EXEC build

Set-Location $CWD

Copy-Item "$FRONT\dist" "$DIST\daily_wallpaper\static" -Recurse -Container

Set-Location $DIST

7z a "daily-wallpaper-windows-amd64-v$VERSION.7z" daily_wallpaper

Set-Location $CWD