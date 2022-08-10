<a href="https://nuxui.org/">
  <h1 align="center">
    <picture>
      <source height="150" media="(prefers-color-scheme: dark)" srcset="https://gitee.com/nuxui/website/raw/master/static/nuxui_logo_text.svg">
      <img alt="NuxUI" height="150" src="https://gitee.com/nuxui/website/raw/master/static/nuxui_logo_text.svg">
    </picture>
  </h1>
</a>

NuxUI is Golang GUI SDK for IOS, Android, macOS, Windows, Linux from a single codebase.

NuxUI is now in developing, the API maybe changed before first stable version.

Any suggestion or good idea post to [discussions](https://github.com/nuxui/nuxui/discussions), let's us make it awesome

## Documentation

* [Install NuxUI](https://nuxui.org/start/install/)

* [NuxUI Documentation](https://nuxui.org/)

* [NuxUI Samples](https://github.com/nuxui/samples)

## Quick Start for desktop
```
git clone https://github.com/nuxui/samples.git
cd github.com/nuxui/samples/widgets
go mod tidy
go build . && ./widgets
```

## Build for mobile
```
# get build tools
go install nuxui.org/nuxui/cmd/nux@latest

# into sample code
cd github.com/nuxui/samples/counter

# ios
nux build -target=iossimulator -bundleid="app.id" -teamid="YOURTEAMID" .
xcrun simctl install booted ./counter.app

#android
nux build -target=android -ldflags="-s -w" . 
adb install -r counter.apk
```

## Screenshot

<img src="https://gitee.com/nuxui/website/raw/master/static/samples/screenshot_widgets.webp" width="400px" >

<img src="https://gitee.com/nuxui/website/raw/master/static/samples/screenshot_ios.webp" width="180px" ><img src="https://gitee.com/nuxui/website/raw/master/static/samples/screenshot_android.webp" width="221px" >
