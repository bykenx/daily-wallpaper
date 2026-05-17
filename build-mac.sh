#!/usr/bin/env sh

APP_NAME="Daily Wallpaper"
BUNDLE_NAME="com.bykenx.daily-wallpaper"
EXEC_NAME="daily_wallpaper"
APP_ICON_SIZES="
16,16x16
32,16x16@2x
32,32x32
64,32x32@2x
128,128x128
256,128x128@2x
256,256x256
512,256x256@2x
512,512x512
1024,512x512@2x
"
NPM_EXEC=""
FRONTEND_DIR="$(pwd)/frontend"
DIST_DIR="$(pwd)/dist"
VERSION=$(cat version.txt)

template=$(cat << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>$EXEC_NAME</string>
	<key>CFBundleIconFile</key>
	<string>AppIcon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>$BUNDLE_NAME</string>
	<key>CFBundleName</key>
	<string>$APP_NAME</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>$VERSION</string>
	<key>CFBundleVersion</key>
	<string>$VERSION</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>LSUIElement</key>
	<true/>
</dict>
</plist>
EOF
)


# 判断依赖
if [ ! -x "$(command -v svg2png)" ];then
  echo "svg2png 未安装，请通过 Homebrew（brew install svg2png）或其他方式安装"
  exit -1
fi
if [ ! -x "$(command -v 7z)" ];then
  echo "7zip 未安装，请通过 Homebrew（brew install p7zip）或其他方式安装"
  exit -1
fi
if [ -x "$(command -v pnpm)" ];then
  NPM_EXEC=pnpm
elif [ -x "$(command -v npm)" ];then
  NPM_EXEC=npm
elif [ -x "$(command -v yarn)" ];then
    NPM_EXEC=yarn
fi
if [ -z "$NPM_EXEC" ];then
  echo "请安装nodejs环境，并确保npm、pnpm、yarn任一包管理器已经安装"
  exit -1
fi

# 清空文件夹
rm -rf ./dist/*

# 构建资源文件
cd "$FRONTEND_DIR"
$NPM_EXEC install
$NPM_EXEC run build
cd ..

build_arch_app() {
    ARCH="$1"
    ARCH_DIR="$DIST_DIR/$ARCH"
    PKG_NAME="$APP_NAME.app"
    PKG_DIR="$ARCH_DIR/$PKG_NAME"
    CONTENT_DIR="$PKG_DIR/Contents"
    RESOURCES_DIR="$CONTENT_DIR/Resources"
    EXEC_DIR="$CONTENT_DIR/MacOS"

    mkdir -p "$RESOURCES_DIR" "$EXEC_DIR"

    # 创建 plist
    echo "$template" > "$CONTENT_DIR/Info.plist"

    # 生成图标文件
    ICONSET="$RESOURCES_DIR/AppIcon.iconset"
    mkdir -p "$ICONSET"
    for PARAMS in $APP_ICON_SIZES; do
        SIZE=$(echo "$PARAMS" | cut -d, -f1)
        LABEL=$(echo "$PARAMS" | cut -d, -f2)
        svg2png -w "$SIZE" -h "$SIZE" buildAssets/AppIcon.svg "$ICONSET/icon_$LABEL.png"
    done
    iconutil -c icns "$ICONSET"
    rm -rf "$ICONSET"

    # 构建指定架构应用，不再通过 lipo 合并成 universal binary
    CGO_ENABLED=1 GOOS=darwin GOARCH="$ARCH" go build -o "$EXEC_DIR/$EXEC_NAME"
    cp -rf "$FRONTEND_DIR/dist" "$RESOURCES_DIR/static"

    # SMAppService (macOS 13+ 登录项) 会校验代码签名；对每个架构包分别做本地 ad-hoc 签名。
    codesign --force --sign - "$EXEC_DIR/$EXEC_NAME"
    codesign --force --sign - "$PKG_DIR"

    (cd "$ARCH_DIR" && 7z a "$DIST_DIR/daily-wallpaper-mac-$ARCH-v${VERSION}.7z" "$PKG_NAME")
}

build_arch_app amd64
build_arch_app arm64