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
