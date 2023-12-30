# Velx Discord Bot

Velx is a Discord bot written in Go (Golang) by KaynHvH using the DiscordGo library. It provides various moderation and fun commands for your Discord server.

## TODO
- Slash commands
- More features
- Better code structuring

## Features

- **Moderation Commands:**
- `velx ban <@user> <reason>` - Bans a user with a specified reason.
- `velx kick <@user> <reason>` - Kicks a user with a specified reason.
- `velx mute <@user> <reason>` - Mutes a user with a specified reason.
- `velx unmute <@user>` - Unmutes a user.
- `velx nick/nickname <@user> <nickname>` - Changes user nickname.
- `velx poll <content>` - Creates a poll with reactions.

- **Fun Commands:**
- `velx dog` - Sends a random dog picture with API.
- `velx answer <question>` - Provides a random answer to a user's question.
- `velx whois <@user>` - Provides information about a Discord user.
- `velx avatar <@user>` - Shows up user's avatar.
- `velx dice` - Rolls the dice.

- **Help Command:**
- `velx help` - Displays a list of available commands and their descriptions.
## Getting Started
1. Clone the repository:
```bash
git clone https://github.com/KaynHvH/velxBOT.git
cd velxBOT
```
2. Create a `.env` file with your Discord bot token and prefix:
```
TOKEN=your_bot_token_here 
PREFIX=velx
```

3. Run the bot:
```bash
make
```

## PREFIX

- Prefix: `velx`
- In the future: slash commands
## Contributing

If you would like to contribute to the project, please follow the [Contributing Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
