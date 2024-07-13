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

3. Run this command in the "main" channel the bot and you are present.
   This sets up all default permissions for the specified channel.
```sh
!managecmd setup #ChannelName
```
### Examples how to use the managecmd command
 
- `!managecmd edit <command> <role> <channels...>`
Edits the specified command to be allowed for the given role in the listed channels.

`!managecmd edit !hello Admin #channel1 #channel2`

- `!managecmd add <command> <role> <channels...>`
Adds a new command with specified role and channels.

`!managecmd add !hello Admin #channel1 #channel2`

- `!managecmd remove <command> <role>`
Removes the specified command for the given role.

`!managecmd remove !hello Admin`

- `!managecmd list <command>`
Lists all permissions for the specified command.

`!managecmd list !hello`



## User Management
Users can be set to a specific role per channel. so they can have diffrent roles in diffrent channels. if no role is set users will be treated as "Everyone" role, this role is a generic one that defaults for everyone that does not have an role set.

### Managing Users

You can manage users by adding or removing them using the following commands in the IRC channels the bot is in:

- **Add User**: `!adduser <nickname> <role> <channel>`
- **Remove User**: `!deluser <nickname> <channel>`

### Roles

The following roles are supported:

- `Owner`: Highest privilege level.
- `Admin`: Administrative privileges.
- `Trusted`: Trusted user.
- `BadBoy`: Restricted user.

## Current Commands

- `!hello` Hello example command available to everyone.
- `!url` Enables/disables URL features (YouTube, Wikipedia, etc.) for admins.
- `!op <user> <channel>` Ops a user in the channel, admin only.
- `!deop <user> <channel>` Deops a user in the channel, admin only.
- `!voice <user> <channel>` Voices a user in the channel, admin only.
- `!devoice <user> <channel>` Devoices a user in the channel, admin only.
- `!kick <user> <channel>` Kicks a user from the channel, admin only.
- `!ban <user> <channel>` Bans a user from the channel, admin only.
- `!unban <user> <channel>` Unbans a user from the channel, admin only.
- `!invite <user> <channel>` Invites a user to the channel, admin only.
- `!topic <new_topic> <channel>` Changes the channel topic, admin only.
- `!join <channel>` Bot joins the specified channel, admin only.
- `!part <channel>` Bot parts from the specified channel, admin only.
- `!shutdown` Shuts down the bot, owner only.
- `!nick <new_nickname>` Changes the bot's nickname, owner only.

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
