# reportportal-dashboards-as-code

A tooling to import and export Report Portal dashboards in YAML

## Build & Install

Prerequisites:

- [go](https://go.dev/)
- [git](https://git-scm.com/)

Clone the repository on your local machine:

```
git clone https://github.com/b1zzu/reportportal-dashboards-as-code.git
```

build the `rpdac` tool for your machine:

```
go build -o rpdac
```

test the new executable:

```
./rpdac help
```

you can now use the `rpdac` as it is or move it in a PATH directory in your machine to install it globally.

## Export a Dashboard

1. Create `.rpdac.toml` configuration file for in the directory from which you will execute the `rpdac` tool like this:

   ```toml
   endpoint = "https://your.report-portal.url"
   token = "your-report-portal-token"
   ```

   > Note: you can also specify a different path for the config file using `-c` flag when executing `rpdac`

1. Identify the Dashboard ID, which you can find in the Dashboard URL if you open the Dashboard you want to export in the ReportPortal UI

1. Export the Dashboard to a local YAML file:

   ```
   rpdac -p the_dashboard_project -d 1 -f ./mydashboard.yaml
   ```

   > Note: in this example, the Dashboard ID we are going to export is `1`

## Import/Create a Dashboard

1. Like for [Export a Dashboard](#export-a-dashboard) you will need to create a `.prdac.toml` config file before proceeding

1. Create a YAML file with the Dashboard definition you want to create

   > Note: use the UI to create a Dashboard and Widgets then export it using `rpdac export ...` to learn the format of a YAML Dashboard, you can then duplicate and/or edit the YAML manually to create multiple Dashboard, or use it as it is too import it in a different Report Portal instance or Project.

1. Create the Dashboard in Report Portal from the YAML definition using the `rpdac create ...` like this:

   ```
   rpdac -p my_project -f ./mydashboard.yaml
   ```
