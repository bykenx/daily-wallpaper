#!/usr/bin/env sh

CWD=$(pwd)
FRONT="$CWD/frontend"
DIST="$CWD/dist"
VERSION=$(cat version.txt)

NPM_EXEC="npm"

install_if_not_exist() {
    if [ -z "$(dpkg -l | grep $1)" ]; then
        sudo apt-get install $1 -y
    fi
}
# 自动安装 build-essential libappindicator3-dev libgtk-3-dev
install_if_not_exist build-essential
install_if_not_exist libgtk-3-dev
install_if_not_exist libappindicator3-dev
install_if_not_exist p7zip-full

# 判断是否安装了 nodejs
if [ -z "$(which node)" ]; then
    echo "未检测到已安装的Node.js，推荐手动安装，继续将自动安装？(y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        sudo apt-get install nodejs -y
    else
        exit 1
    fi
fi


for i in pnpm yarn npm; do
    if [ -n "$(which $i)" ]; then
        NPM_EXEC=$i
        break
    fi
done

rm -rf "$DIST"amd64
go build -o "$DIST/daily_wallpaper/daily_wallpaper_linux_amd64"
cd $FRONT
$NPM_EXEC install
$NPM_EXEC build
cd $CWD
cp -r "$FRONT/dist" "$DIST/daily_wallpaper/static"
cd $DIST
7z a "daily-wallpaper-linux-amd64-v$VERSION.7z" daily_wallpaper
cd $CWD