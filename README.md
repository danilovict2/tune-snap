# Tune Snap 🎵

Tune Snap is a music recognition platform inspired by Shazam’s algorithm. In addition to identifying songs, it integrates Spotify and YouTube APIs to store songs in the library.

<a width="100" href="https://drive.google.com/file/d/19aAHli-5rxUeuvCK7vG2mp-Yw3QKIYLW/view?usp=drive_link">
<img src="https://drive.google.com/file/d/1OQFSVYB7jFFV2xmWHMpAMm6GhCgRyle8/view?usp=drive_link"/>
</a>

<a align="center" href="https://drive.google.com/file/d/19aAHli-5rxUeuvCK7vG2mp-Yw3QKIYLW/view?usp=drive_link">
Demo
</a>

## 🚀 Quick Start

Ensure you have either [Docker](https://www.docker.com/get-started) or the following installed and configured:

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/)
- [ffmpeg](https://ffmpeg.org/)

### Clone the project

```bash
git clone https://github.com/danilovict2/tune-snap.git
cd tune-snap
```

### Set environment variables

```bash
cp .env.example .env
```

### Spotify API

1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) and create a new application.
2. Copy the generated **Client ID** and **Client Secret**.
3. Open your `.env` file and paste these values into the appropriate fields.

### Run with Docker

```bash
docker compose up --build
```

### Run locally

1. Set up a MongoDB instance (locally or using a cloud provider like [MongoDB Atlas](https://www.mongodb.com/atlas)).  
2. Copy the connection string and paste it into the `MONGODB_URI` field in your `.env` file.
3. Install frontend dependencies
```bash
npm install
```

4. Build frontend assets

```bash
npm run build
```

5. Run the Application

```bash
go build -o bin/app
./bin/app
```

### Access the Platform

Open your web browser and go to [http://localhost:8000](http://localhost:8000).

> **Note:** For the best experience, use `localhost` instead of `127.0.0.1`. This ensures embedded YouTube videos work correctly within the platform.


## ✨ Features

- 🎶 Accurate music recognition powered by advanced audio fingerprinting.
- 🔗 Seamless integration with Spotify and YouTube APIs for song discovery and library management.
- ⚙️ Effortless setup using Docker or manual local installation.

## 🤝 Contributing

### Build the project

```bash
npx webpack
go build -o bin/app
```

### Run the project

```bash
./bin/app
```

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
