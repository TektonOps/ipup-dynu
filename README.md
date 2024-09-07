<h2 align="center">
  <br>
  <p align="center"><img width=30% src=".github/img/logo.png"></p>
</h2>


<p align="center">
  <a href="#installation">Install</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#issues">Issues</a> •
</p>

# iPUp Dynu

Dynamic IP updater periodically checks your public IP address and automatically updates your DNS records in **Dynu** if the IP address changes. 


## Installation

### Running the Application

1. **Kubernetes**:
   Install in kubernetes via helm chart

```bash
helm repo add tektonops https://charts.tektonops.com
helm repo update

helm install ipup-dynu tektonops/ipup-dynu \
  --set ipup.config.domainName="yourdomian.com" \
  --set ipup.config.username:="dynuUserName" \
  --set ipup.config.password="dynuPassword"

```

2. **Docker**:
   Install in kubernetes via helm chart

```bash
docker run -d --name ipup \
-e DYNU_DOMAIN_NAME="yourdomina.com" \
-e DYNU_USERNAME="dynuUserName" \
-e DYNU_PASSWORD="dynuPassword" \
tektonops/ipup-dynu:latest

```

## Configuration

The application uses the following environment variables for configuration:

- `DYNU_USERNAME`: Dynu account username.
- `DYNU_PASSWORD`: Dynu account password.
- `DYNU_DOMAIN_NAME`: The domain name you want to update.
- `DYNU_ENABLE_GROUP`: Set to `true` if you want to update a group, `false` otherwise.
- `DYNU_GROUP_NAME`: The group name if you are using a group to update the DNS record.
- `LOG_LEVEL`: Set the logging level (`debug`, `info`, `warn`, `error`). Default (info)
- `DYNU_IPCHECK_INTERVAL`: The time in seconds the application waits between IP checks (defaults to 60 seconds if not set).
- `IPSERVERS_LIST`: Provide your own comma separated list of website which returns your public IP.




#### Config file
- You can also configure the app via config file `config.yaml`. set `USE_CONFIG_FILE=true`
```yaml
dynu:
  domain: "example.com"
  group: "webservers"
  username: "user@example.com"
  password: "154634ed9b592e8a4"
  enableGroup: true
  ipCheckInterval: 60
  ipServersList:
    - https://api4.ipify.org
    - https://ip2location.io/ip
    - https://ident.me

logs:
  logLevel: debug
  enableSource: true
```

## Application Flow

1. The application retrieves the current public IP address from a list of servers.
2. It compares this IP address with the last known IP stored in memory.
3. If the IP has changed, it updates the DNS record using the Dynu API.
4. The application then waits for the duration specified by `DYNU_IPCHECK_INTERVAL` before checking the IP again.

## Issues
If you encounter any problems or have suggestions for improvements, please [open an issue](https://github.com/tektonops/ipup-dynu/issues) on GitHub.

## Contributing

Contributions, issues and feature requests are welcome!<br/>Feel free to check
[issues page](https://github.com/tektonops/ipup-dynu/issues). You can also take a look
at the [contributing guide](https://github.com/tektonops/ipup-dynu/blob/main/CONTRIBUTING.md).