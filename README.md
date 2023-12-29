# Velx Discord Bot

Velx is a Discord bot written in Go (Golang) using the DiscordGo library. It provides various moderation and fun commands for your Discord server.

## Features

- **Moderation Commands:**
- `velx ban <@user> <reason>` - Bans a user with a specified reason.
- `velx kick <@user> <reason>` - Kicks a user with a specified reason.
- `velx mute <@user> <reason>` - Mutes a user with a specified reason.

- **Fun Commands:**
- `velx dog` - Sends a random dog picture.
- `velx answer <question>` - Provides a random answer to a user's question.
- `velx whois <@user>` - Provides information about a Discord user.

- **Help Command:**
- `velx help` - Displays a list of available commands and their descriptions.
## Getting Started
1. Clone the repository:
```bash
git clone https://github.com/yourusername/velx-discord-bot.git
cd velx-discord-bot
```
2. Create a `.env` file with your Discord bot token:
```
TOKEN=your_bot_token_here 
```

3. Run the bot:
```bash
go run main.go
```

## Commands

- Prefix: `velx`

### Moderation Commands

- `velx ban <@user> <reason>` - Bans a user with a specified reason.
- `velx kick <@user> <reason>` - Kicks a user with a specified reason.
- `velx mute <@user> <reason>` - Mutes a user with a specified reason.

### Fun Commands

- `velx dog` - Sends a random dog picture.
- `velx answer <question>` - Provides a random answer to a user's question.
- `velx whois <@user>` - Provides information about a Discord user.

### Help Command

- `velx help` - Displays a list of available commands and their descriptions.

## Contributing

If you would like to contribute to the project, please follow the [Contributing Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.