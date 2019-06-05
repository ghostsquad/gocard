<h1 align="center">
  <br>
  <a href="http://github.com/ghostsquad/gocard"><img src="https://user-images.githubusercontent.com/903488/57989497-42e83f80-7a50-11e9-9452-48d87e93ef10.png" alt="playing card" width="200px" /></a>
  <br>
  GoCard
  <br>
</h1>

<h4 align="center">A minimal tabletop card prototyping app.</h4>

<p align="center">
  <a href="https://saythanks.io/to/ghostsquad">
      <img src="https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg">
  </a>
  <a href="https://www.paypal.me/WMcNamee">
    <img src="https://img.shields.io/badge/$-donate-ff69b4.svg?maxAge=2592000&amp;style=flat">
  </a>
</p>

Successor to [ghostsquad/pycard](https://github.com/ghostsquad/pycard)

## How To Use

```shell
gocard
```

## Quick Install

### Binary (Cross-platform)

Download the appropriate version for your platform from GoCard Releases. Once downloaded, the binary can be run from anywhere. You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.

Ideally, you should install it somewhere in your PATH for easy use. /usr/local/bin is the most probable location.
Homebrew (macOS)

If you are on macOS and using Homebrew, you can install GoCard with the following one-liner:

```shell
brew install gocard
```

For more detailed explanations, read the installation guides that follow for installing on macOS and Windows.

### Chocolatey (Windows)

If you are on a Windows machine and use Chocolatey for package management, you can install GoCard with the following one-liner:

```shell
choco install gocard -confirm
```

### Scoop (Windows)

If you are on a Windows machine and use Scoop for package management, you can install GoCard with the following one-liner:

```shell
scoop install gocard
```

## Development

1. Copy `build.example.env` to `build.env`
2. fill in your Google API Oauth 2.0 Credentials, described here: [Using OAuth 2.0 to Access Google APIs]([Using OAuth 2.0 to Access Google APIs](https://developers.google.com/identity/protocols/OAuth2))
3. Run this

    ```bash
    go mod vendor
    go get github.com/GeertJohan/go.rice/rice

    source ./build.env
    export GOCARD_CLIENT_ID
    export GOCARD_CLIENT_SECRET

    yarn --frozen-lockfile
    rm -rf ./dist
    mkdir -p ./dist/js
    cp ./node_modules/livereload-js/dist/livereload.min.js ./dist/js

    rice embed-go
    go build -ldflags "-X main.ClientID=${GOCARD_CLIENT_ID} -X main.ClientSecret=${GOCARD_CLIENT_SECRET}"
    ```

## Support

<a href="https://www.buymeacoffee.com/50onA1pjc" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="margin-left: 5px; height: auto !important; width: auto !important;" /></a>

<a href="https://www.paypal.me/WMcNamee" target="_blank"><img src="https://user-images.githubusercontent.com/903488/57995914-6ecbeb00-7a79-11e9-9f04-e5c7170a8d8a.png" alt="Paypal Donate" style="height: 65px !important;width: auto !important;" /></a>

_Icons made by [Freepik](https://www.freepik.com) from [www.flaticon.com](http://www.flaticon.com) is licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/)_
