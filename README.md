# Mbot

Mbot is an IRC bot developed by Tommy Mathisen, designed purely for enjoyment and learning.

This project is a fun exploration into the world of IRC, and is currently still under development, with many features yet to be completed.

## Prerequisites

- Go 1.16 or later (recommended to use the latest version)

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

    ````json
    {
      "server": "irc.libera.chat",
      "port": "6697",
      "nick": "ExampleNick",
      "channels": ["#examplechannel"],
      "nick_serv_user": "ExampleNickServUser",
      "nick_serv_pass": "ExampleNickServPass",
      "use_tls": true
    }
    ````

    - `server`: The address of the IRC server you want to connect to (e.g., `irc.libera.chat`).
    - `port`: The port number for the IRC server (usually `6667` for non-TLS, `6697` for TLS).
    - `nick`: The nickname the bot will use on the IRC server.
    - `channels`: A list of channels the bot should join upon connecting (e.g., `["#channel1", "#channel2"]`).
    - `nick_serv_user`: NickServ username, if your IRC server requires NickServ authentication.
    - `nick_serv_pass`: NickServ password, if your IRC server requires NickServ authentication.
    - `use_tls`: A boolean value (`true` or `false`) indicating whether to use TLS for the connection.

4. Run the bot:
    ```sh
    go run main.go
    ```

    The first time the bot starts, if no owner is set in the `users.json` file, the bot will automatically prompt you to set an owner. This process includes the following steps:

    - The bot will display a series of messages indicating that no owner was found and prompt you to enter the owner's nickname and a setup password.
    - The bot will verify the owner's nickname and send a message to the owner to confirm the setup password.
    - You will need to respond to the bot with the correct setup password within 1 minute.
    - Upon successful confirmation, the owner will be added to the `users.json` file with the role of "Owner".
    - If the password is incorrect or no response is received within 1 minute, the bot will shut down.

    Here's what the process looks like:

    ```sh
    =============================== NO OWNER FOUND ===============================
    No owner was found in the users.json file. Please set an owner.
    The bot will shut down if no owner is set within 1 minute after connecting.
    The bot will message the owner to confirm the Setup password.
    ==============================================================================
    >> Please enter the nick of the owner on the network:
    >> Please enter your Setup password:
    ```

    Ensure to follow the prompts and set up the owner properly to avoid the bot shutting down.


5. Set up default permissions for all commands. **This step is crucial to enable all other commands.** If you skip this setup, no commands will work. Use the `!managecmd setup <channel>` command in IRC:
    ```irc
    !managecmd setup #newchannel
    ```

    This command will:

    - Set up default permissions for the specified channel.
    - Create a backup of the current configuration before making changes.
    - Clear existing permissions for the channel.
    - Set new default permissions for the channel.
    - Save and reload the updated command configuration.
    - Re-register commands to reflect updated permissions.

    Example usage:

    ```irc
    !managecmd setup #examplechannel
    ```

    Make sure to run this command in the channel where the bot is present to complete the setup. **Without this step, no other commands will be functional.**

## User Management

### Managing Users

You can manage users by adding or removing them using the following commands in the IRC channels the bot is in:

- **Add User**: `!adduser <nickname> <role>`
- **Remove User**: `!deluser <nickname>`

### Roles

The following roles are supported:

- `Owner`: Highest privilege level.
- `Admin`: Administrative privileges.
- `Trusted`: Trusted user.
- `BadBoy`: Restricted user.

### Example `users.json`

The users are stored in a `users.json` file located in the `data` directory. Below is an example structure:


```json
{
    "~jane@irc/example/com/jane": {
      "hostmask": "~jane@irc/example/com/jane",
      "roles": {
        "*": "Owner"
      }
    },
    "~bob@irc/example/com/bob": {
      "hostmask": "~bob@irc/example/com/bob",
      "roles": {
        "#general": "Trusted",
        "#support": "Trusted"
      }
    },
    "eve@irc/example/com/eve": {
      "hostmask": "eve@irc/example/com/eve",
      "roles": {
        "#general": "BadBoy",
        "#random": "BadBoy"
      }
    },
    "charlie@irc/example/com/charlie": {
      "hostmask": "charlie@irc/example/com/charlie",
      "roles": {
        "#admin": "Admin",
        "#dev": "Trusted"
      }
    }
  }
```

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
## Credits

- This bot uses the [ErgoChat IRC library](https://github.com/ergochat/ergo).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.
