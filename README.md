# Discord channel recorder

## Before running the bot please configure the following fields with your own configuration settings:

```bash
    export BOT_TOKEN="your_bot_token"
    export TIME_LIMIT=5
    export S3_REGION="eu-north-1"
    export S3_ACCESS_KEY="your_access_key"
    export S3_SECRET_KEY="your_secret_key"
    export S3_BUCKET_NAME="your_bucket_name"
```

## To run the bot

```bash
    go run ./cmd/
```

## NOTE: please configure your Amazon S3 by giving permission to upload files to your bucket

## Methodology

1. I used clean architecture for this project, since the bot also operates as backend service, where we have handler, service and repository layers. However, I need to mention that I did not properly implemented the service layer by putting the service functionality in the same folder as hadndler part.
2. Used libraries:
   - Discordgo for interacting with Discord
   - Pion for recording audio from voice channel and recording it to a file
   - Amazon S3 SDK for uploading files and getting their URLs
   - Logrus for logging
