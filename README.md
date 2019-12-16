# `fvm`

Flutter Version Management: A simple cli to manage Flutter SDK versions.

**Features:**

* Configure Flutter SDK version per project or globally
* Ability to install and cache multiple Flutter SDK Versions
* Easily switch between Flutter channels & versions

## Version Management

This tool allows you to manage multiple channels and releases, and caches these versions locally, so you don't have to wait for a full setup every time you want to switch versions.

Also, it allows you to grab versions by a specific release, i.e. 1.2.0. In case you have projects in different Flutter SDK versions and do not want to upgrade.

## Usage

To Install fvm:

```bash
> go get -u -v github.com/befovy/fvm
```

And then, for information on each command:

```bash
> fvm help
```

### FVM_HOME

fvm use environment variable `FVM_HOME` as a working path.  
The installed cache and config file are all stored in this path.

If no `FVM_HOME` in environment variable, fvm will use the default value returned by `os.UserConfigDir()` append `fvm`.  
On Mac OS, default FVM_HOME is $HOME/Library/Application\ Support/fvm

### Install a SDK Version

FVM gives you the ability to install many Flutter **releases** or **channels**.

```bash
> fvm install <version>
```

Use `master` to install the Master channel and `v1.8.0` to install the release.

### Use a SDK Version

You can use the installed Flutter SDK versions for your computer user account globally. To do that just:

```bash
> fvm use <version>
```

Also, you can use different Flutter SDK versions per project. To do that you have to go into the root of the project and:

```bash
> fvm use <version> --locol
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

### Running Flutter SDK

#### Call Global SDK 

After add fvm global path to your system environment variable, `flutter` command is usable every where.

Your will get tip when you run `fvm use <version>` firstly.



#### Proxy Commands

Flutter command within `fvm` proxies all calls to the CLI just changing the SDK to be the local one.

```bash
> fvm flutter run
```

This will run `flutter run` command using the local project SDK. If no FVM config is found in the project. FMV will recursively try for a version in a parent directory.

If FVM config is still not found, this will run `flutter run` command using the global Flutter SDK. 



#### Call Local SDK Directly

FVM creates a symbolic link within your project called **.fvmbin/flutter** which links to the installed version of the SDK.


Add `$(pwd)/fvmbin` to your PATH, or
```bash
> ./fvmbin/flutter run
```

This will run `flutter run` command using the local project SDK.

As an example calling `fvm flutter run` is the equivalent of calling `flutter run` using the local project SDK.



## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details


