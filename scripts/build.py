#!/usr/bin/env python3
"""Cross-platform build script for daily-wallpaper."""

from __future__ import annotations

import argparse
import os
import platform
import shutil
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
ASSETS_DIR = ROOT / "assets"
FRONTEND_DIR = ROOT / "frontend"
DIST_DIR = ROOT / "dist"
APP_NAME = "Daily Wallpaper"
BUNDLE_NAME = "com.bykenx.daily-wallpaper"
EXEC_NAME = "daily_wallpaper"
GO_TARGET = "./cmd/daily-wallpaper"
APP_ICON_SIZES = [
    ("16", "16x16"),
    ("32", "16x16@2x"),
    ("32", "32x32"),
    ("64", "32x32@2x"),
    ("128", "128x128"),
    ("256", "128x128@2x"),
    ("256", "256x256"),
    ("512", "256x256@2x"),
    ("512", "512x512"),
    ("1024", "512x512@2x"),
]


def run(command: list[str], *, cwd: Path = ROOT, env: dict[str, str] | None = None) -> None:
    print("+", " ".join(command))
    subprocess.run(command, cwd=cwd, env=env, check=True)


def require(command: str) -> str:
    path = shutil.which(command)
    if path is None:
        raise SystemExit(f"Missing required command: {command}")
    return path


def version() -> str:
    return (ROOT / "version.txt").read_text().strip()


def npm_command() -> str:
    for command in ("pnpm", "npm", "yarn"):
        path = shutil.which(command)
        if path is not None:
            return path
    raise SystemExit("Missing Node.js package manager: pnpm, npm, or yarn")


def clean_dist() -> None:
    if DIST_DIR.exists():
        shutil.rmtree(DIST_DIR)
    DIST_DIR.mkdir(parents=True)


def copy_static(target: Path) -> None:
    source = FRONTEND_DIR / "dist"
    if not source.exists():
        raise SystemExit(f"Frontend build output not found: {source}")
    if target.exists():
        shutil.rmtree(target)
    shutil.copytree(source, target)


def build_frontend() -> None:
    npm = npm_command()
    run([npm, "install"], cwd=FRONTEND_DIR)
    run([npm, "run", "build"], cwd=FRONTEND_DIR)


def go_build(
    output: Path,
    *,
    goos: str,
    goarch: str,
    cgo_enabled: str | None = None,
    ldflags: str | None = None,
) -> None:
    output.parent.mkdir(parents=True, exist_ok=True)
    env = os.environ.copy()
    env.update({"GOOS": goos, "GOARCH": goarch})
    if cgo_enabled is not None:
        env["CGO_ENABLED"] = cgo_enabled
    command = ["go", "build", "-o", str(output)]
    if ldflags:
        command.append(f"-ldflags={ldflags}")
    command.append(GO_TARGET)
    run(command, env=env)


def archive(source: Path, output: Path) -> None:
    require("7z")
    if output.exists():
        output.unlink()
    run(["7z", "a", str(output), source.name], cwd=source.parent)


def plist_template(app_version: str) -> str:
    return f"""<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
\t<key>CFBundleExecutable</key>
\t<string>{EXEC_NAME}</string>
\t<key>CFBundleIconFile</key>
\t<string>AppIcon.icns</string>
\t<key>CFBundleIdentifier</key>
\t<string>{BUNDLE_NAME}</string>
\t<key>CFBundleName</key>
\t<string>{APP_NAME}</string>
\t<key>CFBundlePackageType</key>
\t<string>APPL</string>
\t<key>CFBundleShortVersionString</key>
\t<string>{app_version}</string>
\t<key>CFBundleVersion</key>
\t<string>{app_version}</string>
\t<key>NSHighResolutionCapable</key>
\t<true/>
\t<key>LSUIElement</key>
\t<true/>
</dict>
</plist>
"""


def build_macos(arches: list[str]) -> None:
    require("svg2png")
    require("iconutil")
    require("codesign")
    build_frontend()
    clean_dist()
    app_version = version()
    for arch in arches:
        arch_dir = DIST_DIR / arch
        package_dir = arch_dir / f"{APP_NAME}.app"
        contents_dir = package_dir / "Contents"
        resources_dir = contents_dir / "Resources"
        executable_dir = contents_dir / "MacOS"
        resources_dir.mkdir(parents=True, exist_ok=True)
        executable_dir.mkdir(parents=True, exist_ok=True)

        (contents_dir / "Info.plist").write_text(plist_template(app_version))

        iconset = resources_dir / "AppIcon.iconset"
        iconset.mkdir()
        for size, label in APP_ICON_SIZES:
            run(
                [
                    "svg2png",
                    "-w",
                    size,
                    "-h",
                    size,
                    str(ASSETS_DIR / "AppIcon.svg"),
                    str(iconset / f"icon_{label}.png"),
                ]
            )
        run(["iconutil", "-c", "icns", str(iconset)])
        shutil.rmtree(iconset)

        executable = executable_dir / EXEC_NAME
        go_build(executable, goos="darwin", goarch=arch, cgo_enabled="1")
        copy_static(resources_dir / "static")
        run(["codesign", "--force", "--sign", "-", str(executable)])
        run(["codesign", "--force", "--sign", "-", str(package_dir)])
        archive(package_dir, DIST_DIR / f"daily-wallpaper-mac-{arch}-v{app_version}.7z")


def build_linux(arch: str) -> None:
    build_frontend()
    clean_dist()
    app_version = version()
    package_dir = DIST_DIR / "daily_wallpaper"
    go_build(package_dir / f"daily_wallpaper_linux_{arch}", goos="linux", goarch=arch)
    copy_static(package_dir / "static")
    archive(package_dir, DIST_DIR / f"daily-wallpaper-linux-{arch}-v{app_version}.7z")


def build_windows(arch: str) -> None:
    build_frontend()
    clean_dist()
    app_version = version()
    package_dir = DIST_DIR / "daily_wallpaper"
    go_build(
        package_dir / "daily_wallpaper.exe",
        goos="windows",
        goarch=arch,
        ldflags="-H windowsgui",
    )
    copy_static(package_dir / "static")
    archive(package_dir, DIST_DIR / f"daily-wallpaper-windows-{arch}-v{app_version}.7z")


def current_platform() -> str:
    name = platform.system().lower()
    if name == "darwin":
        return "darwin"
    if name == "linux":
        return "linux"
    if name == "windows":
        return "windows"
    raise SystemExit(f"Unsupported platform: {platform.system()}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Build daily-wallpaper packages.")
    parser.add_argument(
        "--arch",
        choices=("amd64", "arm64", "all"),
        default="all",
        help="target architecture; macOS defaults to all, other platforms use amd64 when all is selected",
    )
    return parser.parse_args()


def main() -> None:
    require("go")
    args = parse_args()
    target = current_platform()
    if target == "darwin":
        arches = ["amd64", "arm64"] if args.arch == "all" else [args.arch]
        build_macos(arches)
    elif target == "linux":
        build_linux("amd64" if args.arch == "all" else args.arch)
    elif target == "windows":
        build_windows("amd64" if args.arch == "all" else args.arch)
    else:
        raise SystemExit(f"Unsupported platform: {target}")


if __name__ == "__main__":
    try:
        main()
    except subprocess.CalledProcessError as exc:
        sys.exit(exc.returncode)
