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
      "url_features": {
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

## API Keys

For certain features, you will need to have API keys for YouTube, IMDb, and VirusTotal. These keys should be stored in a `.env` file or exported as environment variables.

## Obtaining API Keys

- **YouTube API Key**: You can get a YouTube API key by creating a project on the [Google Developers Console](https://console.developers.google.com/). Enable the YouTube Data API v3 and generate an API key.

- **OMDb API Key**: To obtain an OMDb API key, register for a free or paid account on the [OMDb website](https://www.omdbapi.com/apikey.aspx). Once registered, you will receive an API key via email.

- **VirusTotal API Key**: You can get a VirusTotal API key by signing up for a free account on the [VirusTotal website](https://www.virustotal.com/). After logging in, navigate to your profile and generate an API key.

### Example `.env` file

Create a `.env` file in the root directory of the project with the following content:

```sh
YOUTUBE_API_KEY=your_youtube_api_key
OMDb_API_KEY=your_omdb_api_key
VIRUSTOTAL_API_KEY=your_virustotal_api_key
```

### Setting Environment Variables

Alternatively, you can set the environment variables directly.

#### On Windows:

```sh
setx YOUTUBE_API_KEY "your_youtube_api_key"
setx OMDb_API_KEY "your_omdb_api_key"
setx VIRUSTOTAL_API_KEY "your_virustotal_api_key"
```

#### On Linux or macOS:

```bash
export YOUTUBE_API_KEY="your_youtube_api_key"
export OMDb_API_KEY="your_omdb_api_key"
export VIRUSTOTAL_API_KEY="your_virustotal_api_key"
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.
