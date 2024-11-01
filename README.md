# UniFi Captive Portal

The [UniFi](https://ui.com/) external captive portal with ChatGPT.

## Demo

![Chat Demo](/doc/demo.png "Demo")

## Configuration

### How to enbable the external portal
Follow up this document: [UniFi Hotspot Portal and Guest WiFi](https://help.ui.com/hc/en-us/articles/115000166827-UniFi-Hotspot-Portal-and-Guest-WiFi)

### Portal Env variables

| Variable        | Description                                                                               | Default                      |
|-----------------|-------------------------------------------------------------------------------------------|------------------------------|
| `UnifiURL`      | Full URL of your UniFi Controller. Be sure to include the port it is running on           | `https://unifi:8443`         |
| `UnifiUsername` | Username of the user to make API calls with. It is recommended to use a dedicated user    | `admin`                      |
| `UnifiPassword` | Password for user defined above                                                           | `admin`                      |
| `UnifiSite`     | The name of the site the APs/Users reside in                                              | `default`                    |
| `Title`         | Title used in HTML pages as well as headings. Usually you will put your company name here | `Captive Portal`             |
| `RedirectUrl`   | URL to redirect users to if they do not provide one to the controller                     | `https://captive.apple.com/` |
| `Minutes`       | Amount of time to register user for                                                       | `600`                        |

### OpenAi Env variables

| Variable      | Description                                                                     | Default                  |
|---------------|---------------------------------------------------------------------------------|--------------------------|
| `ApiKey`      | Your OpenAI API key                                                             | `None`                   |
| `Model`       | The model to use for generating responses                                       | `gpt-3.5-turbo`          |
| `Domain`      | The domain to use for the model                                                 | `https://api.openai.com` |
| `Temperature` | The randomness of the responses. 0.0 is deterministic, 1.0 is completely random | `0.7`                    |

## Build Binary

```shell
$ make build
```

## Build Docker

```shell
$ make docker
```

## Running the Application

To run the application, you need to set the required environment variables and then execute the binary or run the Docker
container.

### Using Binary

1. Set environment variables:
    ```shell
    export UnifiURL=https://unifi:8443
    export UnifiUsername=admin
    export UnifiPassword=admin
    export UnifiSite=default
    export Title="Captive Portal"
    export RedirectUrl=https://captive.apple.com/
    export Minutes=600
    export ApiKey=your_openai_api_key
    export Model=gpt-3.5-turbo
    export Domain=https://api.openai.com
    export Temperature=0.7
    ```

2. Run the binary:
    ```shell
    ./unifi-captive-portal
    ```

### Using Docker

1. Run the Docker container:
    ```shell
    docker run -e UnifiURL=https://unifi:8443 \
               -e UnifiUsername=admin \
               -e UnifiPassword=admin \
               -e UnifiSite=default \
               -e Title="Captive Portal" \
               -e RedirectUrl=https://captive.apple.com/ \
               -e Minutes=600 \
               -e ApiKey=your_openai_api_key \
               -e Model=gpt-3.5-turbo \
               -e Domain=https://api.openai.com \
               -e Temperature=0.7 \
               chinaboard/unifi-captive-portal
    ```

1. Run the Docker compose:
    ```shell
    $ docker compose up -d
    ```
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
