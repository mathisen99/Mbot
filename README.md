# Mbot

IRC bot created by Tommy Mathisen for fun.
The whole project is just for learning and having fun on IRC.

## Prerequisites

- Go 1.16 or later

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/mathisen99/Mbot.git
    cd Mbot
    ```

2. Rename the example configuration file:
    ```sh
    mv data/config_example.json data/config.json
    ```

3. Edit `data/config.json` with your preferred text editor and fill in the configuration details:
    ```json
    {
      "server": "irc.libera.chat",
      "port": "6697",
      "nick": "BotNickname",
      "channels": ["#channel1", "#channel2"],
      "nick_serv_user": "NickServUser",
      "nick_serv_pass": "NickServPass",
      "use_tls": true,
      "features": {
        "enable_youtube_check": true,
        "enable_wikipedia_check": true,
        "enable_github_check": true,
        "enable_imdb_check": true,
        "enable_virus_total_check": true
      }
    }
    ```

    - `server`: The address of the IRC server you want to connect to (e.g., `irc.libera.chat`).
    - `port`: The port number for the IRC server (usually `6667` for non-TLS, `6697` for TLS).
    - `nick`: The nickname the bot will use on the IRC server.
    - `channels`: A list of channels the bot should join upon connecting (e.g., `["#channel1", "#channel2"]`).
    - `nick_serv_user`: NickServ username, if your IRC server requires NickServ authentication.
    - `nick_serv_pass`: NickServ password, if your IRC server requires NickServ authentication.
    - `use_tls`: A boolean value (`true` or `false`) indicating whether to use TLS for the connection.
    - `url_features`: A block containing boolean values to enable or disable specific URL features.
        - `enable_youtube_check`: Enable or disable YouTube link handling.
        - `enable_wikipedia_check`: Enable or disable Wikipedia link handling.
        - `enable_github_check`: Enable or disable GitHub link handling.
        - `enable_imdb_check`: Enable or disable IMDb link handling.
        - `enable_virus_total_check`: Enable or disable VirusTotal link checking.

## Running the Bot

1. Build and run the bot:
    ```sh
    go run .
    ```

    Alternatively, you can build the bot and run the executable:
    ```sh
    go build -o irc-bot
    ./irc-bot
    ```
