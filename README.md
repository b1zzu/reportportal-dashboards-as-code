# reportportal-dashboards-as-code

A tooling to import and export Report Portal dashboards in YAML

## Installation

### Linux

#### Manual installation

1. Navigate to the [Releases](https://github.com/b1zzu/reportportal-dashboards-as-code/releases) page and download the `rpdac_x.y.z_linux_arch.tar.gz` archive, where `x.y.z` is the version you want to install and `arch` is the CPU architecture of your PC (for Intel and AMD CPUs use `amd64`)

2. Extract the `rpdac` from the archive:
   ```
   tar xfv rpdac_x.y.z_linux_arch.tar.gz
   ```

3. Install `rpdac` in your `$PATH` (we suggest installing it in `~/.local/bin`)
   ```
   mv rpdac ~/.local/bin
   rpdac --help
   ```

#### Fedora/CentOS/RHEL/OpenSUSE

1. Navigate to the [Releases](https://github.com/b1zzu/reportportal-dashboards-as-code/releases) page and download the `rpdac_x.y.z_linux_arch.rpm` RPM package, where `x.y.z` is the version you want to install and `arch` is the CPU architecture of your PC (for Intel and AMD CPUs use `amd64`)

2. Install the RPM package:
   ```
   sudo rpm -i ./rpdac_x.y.z_linux_arch.rpm
   rpdac --help
   ```

#### Ubuntu/Debian/ElementaryOS

1. Navigate to the [Releases](https://github.com/b1zzu/reportportal-dashboards-as-code/releases) page and download the `rpdac_x.y.z_linux_arch.deb` DEB package, where `x.y.z` is the version you want to install and `arch` is the CPU architecture of your PC (for Intel and AMD CPUs use `amd64`)

2. Install the RPM package:
   ```
   sudo apt install ./rpdac_x.y.z_linux_arch.deb
   rpdac --help
   ```

### MacOS

#### Manual installation

1. Navigate to the [Releases](https://github.com/b1zzu/reportportal-dashboards-as-code/releases) page and download the `rpdac_x.y.z_darwin_arch.tar.gz` archive, where `x.y.z` is the version you want to install and `arch` is the CPU architecture of your Mac (for Intel CPUs use `amd64`, for Apple M1 use `arm64`)

2. Extract the `rpdac` from the archive:
   ```
   tar xfv rpdac_x.y.z_darwin_arch.tar.gz
   ```

3. Install `rpdac` in your `$PATH` (we suggest installing it in `$HOME/bin`)
   ```
   mv rpdac $HOME/bin
   rpdac --help
   ```

## Getting Started

This is a short step-by-step tutorial on how to use `rpdac`.

**Requirements:**

* `rpdac` is [installed](#Installation) on your PC
* You have a ReportPortal instance and a Access token

1. Navigate to your ReportPortal instance and in your project manually create a dashboard with as many widgets as you want and name it `My Dashboard`.

2. Open a terminal and export your ReportPortal endpoint, Access token, and Project name as environment variables.
   ```
   export RPDAC_ENDPOINT="https://example.com"
   export RPDAC_TOKEN="a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1"
   export RPDAC_PROJECT="my_project"
   ```

3. Run the `export dashboard` command to export the `My Dashboard` locally.
   ```
   rpdac export dashboard --name 'My Dashboard' -f my-dashboard.yaml
   ```

4. Now edit the `description` of the dashboard in the `my-dashboard.yaml` file.

5. Run the `apply` command to update the dashboard in ReportPortal.
   ```
   rpdac apply -f my-dashboard.yaml
   ```

6. Navigate to ReportPortal and verify that the description of `My Dashboard` has updated correctly.



1. Install `rpdac` on your machine

## Configuration

After installing `rpdac` you need to configure the ReportPortal endpoint and token. 

The **Endpoint** is the URL of your ReportPortal instance (example: `https://example.com`)

The **Access token** can be found in the *Profile* page (example: `a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1`)

The endpoint and token can be configured as Environment Variables, Config File, or as Flags before each Command.

### As Environment Variables

| Name | Description |
| ---- | ----------- |
| RPDAC_ENDPOINT | The ReportPortal endpoint URL | 
| RPDAC_TOKEN    | The Access token to authenticate against ReportPortal |

Example:
```
export RPDAC_ENDPOINT="https://example.com"
export RPDAC_TOKEN="a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1"
```

### As Config File

The `.rpdac.toml` config file can be created in the current directory with the ReportPortal endpoint and token:

File example:
```toml
endpoint = "https://example.com"
token = "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1"
```

By default `rpdac` will look for the `.rpdac.toml` config file in the current directory but with the `--config` flag or the `RPDAC_CONFIG` ENV a different path can be specified.

### As Flags before each Command

The `--endpoint` and `--token` flags can also be specified before each command to configure the ReportPortal endpoint and token.

Example:
```
$ rpdac --endpoint "https://example.com" --token "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1" export dashboard [...]
```

## Commands

### Export a Dashboard

Dashboards can be exported in YAML using the `export dashboard` command. Exporting a dashboard is useful to copy it (if you create it again with a different name using the `create` command), save it as code (to be then re-created in the same instance if it gets deleted or in a new instance), migrate it to a new instance.

Example:
```
$ rpdac export dashboard -p my_project --name 'My Dashboard Name' -f my-dashboard.yaml
```

### Export a Filter

Like for Dashboards, Filters can be exported in YAML using the `export filter` command. 

Example:
```
$ rpdac export filter -p my_project --name 'My Filter Name' -f my-filter.yaml
```

### Import/Create a Dashboard

If you already have a Dashboard definition in YAML or you have exported a Dashboard in YAML you can create it in a new ReportPortal instance or in the same if it got deleted using the `create` command.

Example:
```
$ rpdac create -p my_project  -f my-dashboard.yaml
```

> Note: the `create` command automatically detects that the `.yaml` file is a Dashboard

### Import/Create a Filter

Like for Dashboards, Filters can be created from a YAML definition using the `create` command.

Example:
```
$ rpdac create -p my_project  -f my-filter.yaml
```

> Note: the `create` command automatically detect that the `.yaml` file is a Filter

### Apply all Dashboards and Filters from a folder

Using the `apply` command is possible to create or update a single Dashboard or Filter but also an entire directory containing multiple Dashboards and/or Filters.

For example, the current (`$PWD`) directory contains the `my-dashboard.yaml`, `my-filter-01.yaml` and `my-filter-02.yaml`. `my-dashboard.yaml` already exists in ReportPortal but it needs to be updated, `my-filter-01.yaml` doesn't exist in ReportPortal and `my-filter-02.yaml` already exists, and it doesn't need to be updated. When we apply the current directory, the `my-dashboard.yaml` will be updated, the `my-filter-01.yaml` will be created, and the `my-filter-02.yaml` will be ignored.

```
$ rpdac apply -p my_project -f . -r
0000/00/00 00:00:00 Dashboard with name 'My Dashboard' updated in project 'my_project'
0000/00/00 00:00:00 Filter with name 'My Filter 01' created in project 'my_project'
0000/00/00 00:00:00 Skip apply Filter with name 'My Filter 02' in project 'my_project'
```

> Note: If you apply a directory with multiple Dashboards and then you delete one of the Dashboards and apply again the dashboard will not be deleted from ReportPortal, same for filters.

> Note: The apply command will only update a dashboard if it match the name, so if you rename a dashboard in the yaml and apply again it will create a new dashboard in ReportPortal instead of renaming it, same for filters.
