# JAMZ

JAMZ is a terminal tool to manage your spotify playback need without leaving your terminal, perfect for choosing those productivity boosting songs and playlists without disrupting your workflow.

# Getting started

What you need to run the program

- [Go 1.16 or higher](https://go.dev/dl/)
- [Spotify Credentials](https://developer.spotify.com/dashboard) associated with a valid Spotify Premium account

If building locally set up your Spotify credentials in a `credentials.yaml` file as shown in the sample file provided.
Alternatively, you can set up environment variables with your `SPOTIFY_SECRET` and `SPOTIFY_ID` to authenticate the app.
You will need to click the provided link to authenticate the app from the Spotify dashboard to get started.

## How to run

### Locally

1. Clone this git repo
   `git clone https://github.com/Smelton01/jamz`
2. Install/Build the binary
   ```
   cd jamz
   go build or go install
   ```
3. Run the program
   ```
   jamz [command]
   ```

### Docker

1. Build image

```
Docker build --tag jamz .
```

2. Run the container

```
Docker run -d --name jamz
```

## Usage

The executable comes with a few CLI commnds for quick access and a TUI for browsing through your music

### CLI

Run `jamz [command]`

#### Supported commands

- `play` resume playback on your currently active device
- `pause` pause playback
- `next` skip to the next track in your queue
- `prev` skip back to the previous track

### TUI

#### TODO

## Contributing

Contributions are very welcome so feel free to open an issue or submit a PR if you would like to contribute to the project.

## License

[MIT License](https://github.com/Smelton01/jamz/blob/master/LICENSE)

## Credits

Jamz is built using:

- [BubbleTea](https://github.com/charmbracelet/bubbletea) for the UI
- [Spotify](https://github.com/zmb3/spotify) API wrapper
