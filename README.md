# Marketing Bot

This repository contains a bot that interacts with messaging apps like Facebook Messenger. Currently, it is designed to work with Facebook Messenger but has the potential to extend to other messaging platforms.

## Getting Started

To use this bot, you need to create an app on the Facebook Developers platform. Follow these steps:

1. Go to [Facebook for Developers](https://developers.facebook.com/).
2. Create a new app and configure it for Messenger.
3. Retrieve your Facebook Page Access Token and Verify Token.

## Endpoints

### Facebook Endpoints

- `GET /api/v1/facebook/webhook`: Verifies the token.
- `POST /api/v1/facebook/webhook`: Handles incoming messages.
- `POST /api/v1/facebook/upload`: Uploads an image to Facebook. It expects a URL like this: `https://m.media-amazon.com/images/I/2128q5aAVQL.png`.

### Dashboard Endpoints

- `GET /api/v1/dash/clicks`: Returns the count of all clicks on a specific button per user. It requires three query parameters: `question`, `platform`, and `option`. **Note: Use the following mock flow parameters for testing:** `/api/v1/dash/clicks?question=welcome_demo&platform=FB&option=yes`
- `GET /api/v1/dash/revels`: Returns all coupon reveals. It requires two query parameters: `platform` and `code`. **Note: Use the following mock flow parameters for testing:** `/api/v1/dash/revels?platform=FB&code=10FF`

## Dependencies

All dependencies are mocked except for the `GraphAPI`, which interacts with Messenger. Dash info is saved in a memory repository, and the message flow is mocked. To replace a message, go to `fb/service/message_flow.go` on line 332 and change `image_id` to the uploaded image id.

## Usage

To run the bot locally, use the following Makefile commands:

- `make build`: Build the app.
- `make up`: Run the app.
- `make down`: Stop the app.

Ensure you have the following environment variables in your `.env` file:

- `FB_PAGE_ACCESS_TOKEN`: Your Facebook Page Access Token.
- `FB_VERIFY_TOKEN`: Your Facebook Verify Token.
- `PORT`: The port on which the app will run.