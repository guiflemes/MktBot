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


### Flow Endpoints

- `POST /api/v1/flow`: Used to create a flow. This is the mock version. Example payload:
  ```
  {
		"name": "SampleFlow",
		"key": "sample_flow_key",
		"cards": {
		  "buttonCard1": {
			"id": "buttonCard1",
			"type": "button",
			"initial": true,
			"expected_msg": "",
			"template": {
			  "key": "welcome_demo",
			  "text": "Welcome to the demo promotional flow {name}! Are you interested in our coupon",
			  "options": [
				{"text": "Yes! Show me coupon", "target_card_id": "couponCard", "key": "yes"},
				{"text": "No, thanks", "target_card_id": "imageCard", "key": "no"}
			  ]
			}
		  },
		  "imageCard": {
			"id": "imageCard",
			"type": "image",
			"initial": false,
			"expected_msg": "",
			"template": {
			  "image_url": "",
			  "image_id": "1746931029114090"
			}
		  },
		  "couponCard": {
			"id": "couponCard",
			"type": "coupon",
			"initial": false,
			"expected_msg": "",
			"template": {
			  "title": "here is our unqiue promotinal coupon",
			  "subtitle": "10% off limit 1 per customer",
			  "code": "10FF",
			  "key": "10FF"
			}
		  }
		},
		"relationships": [
		  {
			"source_card_id": "buttonCard1",
			"target_card_id": "imageCard",
			"relationship_type": "button_to_image",
			"additional_details": "Option 1 selected"
		  },
		  {
			"source_card_id": "buttonCard1",
			"target_card_id": "couponCard",
			"relationship_type": "button_to_coupon",
			"additional_details": "Option 2 selected"
		  }
		]
	  }
  ```

  `GET /api/v1/flow/:key`: Used to get a flow

## Dependencies

All dependencies are mocked except for the `GraphAPI`, which interacts with Messenger, all repositories are in memory.

## Usage


Ensure you have the following environment variables in your `.env` file:

- `FB_PAGE_ACCESS_TOKEN`: Your Facebook Page Access Token.
- `FB_VERIFY_TOKEN`: Your Facebook Verify Token.
- `PORT`: The port on which the app will run.

* To start, please upload an image. Then, create a flow by replacing the image ID obtained previously in the flow payload. For testing purposes, create the flow exactly like the mock to avoid errors. Only sample_flow_key can be used as the flow key because it saves only one flow in memory and retrieves it based on sample_flow_key.
* Be aware that all the data is in memory so if the server goes down or changes the data is lost and you have to upload the image and create the flow again


To run the bot locally, use the following Makefile commands:

- `make build`: Build the app.
- `make up`: Run the app.
- `make down`: Stop the app.

