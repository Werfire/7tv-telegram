# 7tv-telegram
A Telegram bot that gives you 7TV Emotes

## Using the bot:

Start by typing `@s7ntvbot` in any chat. You can specify optional number for limit as prefix (50 max, default 20) and then optional @CHANNEL_NAME to search emotes from some channel with its aliases. You can also make search exact match and case-sensitive by quoting emote name. Non-ASCII (e.g. cyrillic) emote names will also work.

Example with all the features: `@s7ntvbot 67 @melharucos "melHi"`

## Development

Uses:

* [tucnak/telebot](https://github.com/tucnak/telebot)
* [hasura/go-graphql-client](https://github.com/hasura/go-graphql-client)
