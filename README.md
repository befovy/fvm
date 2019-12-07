# `fvm`

Flutter Version Management: A simple cli to manage Flutter SDK versions.

**Go port of [leoafarias/fvm](https://github.com/leoafarias/fvm)**

**Features:**

* Configure Flutter SDK version per project
* Ability to install and cache multiple Flutter SDK Versions
* Easily switch between Flutter channels & versions
* Per project Flutter SDK upgrade

## Version Management

This tool allows you to manage multiple channels and releases, and caches these versions locally, so you don't have to wait for a full setup every time you want to switch versions.

Also, it allows you to grab versions by a specific release, i.e. 1.2.0. In case you have projects in different Flutter SDK versions and do not want to upgrade.

## Usage

To Install:

```bash
> go get -u -v github.com/befovy/fvm
```

And then, for information on each command:

```bash
> fvm help
```

### Install a SDK Version

FVM gives you the ability to install many Flutter **releases** or **channels**.

```bash
> fvm install <version>
```

Version - use `master` to install the Master channel and `v1.8.0` to install the release.

### Use a SDK Version

You can use different Flutter SDK versions per project. To do that you have to go into the root of the project and:

```bash
> fvm use <version>
```

### Remove a SDK Version

Using the remove command will uninstall the SDK version locally. This will impact any projects that depend on that version of the SDK.

```bash
> fvm remove <version>
```

### List Installed Versions

List all the versions that are installed on your machine.

```bash
> fvm list
```

### Change FVM Cache Directory

There are some configurations that allows for added flexibility on FVM.

```bash
fvm config --cache-path <path-to-use>
```

### List Config Options

Returns list of all stored options in the config file.

```bash
fvm config --ls
```

### Running Flutter SDK

There are a couple of ways you can interact with the SDK setup in your project. 

#### Proxy Commands

Flutter command within `fvm` proxies all calls to the CLI just changing the SDK to be the local one.

```bash
> fvm flutter run
```

This will run `flutter run` command using the local project SDK. If no FVM config is found in the project. FMV will recursively try for a version in a parent directory.

#### Call Local SDK Directly

FVM creates a symbolic link within your project called **fvm** which links to the installed version of the SDK.

```bash
> ./fvm run
```

This will run `flutter run` command using the local project SDK.

As an example calling `fvm flutter run` is the equivalent of calling `flutter run` using the local project SDK.

### Configure Your IDE

#### VSCode

Add the following to your settings.json

```json

"dart.flutterSdkPaths": [
    "fvm"
]
```

## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
