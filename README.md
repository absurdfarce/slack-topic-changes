# slack-topic-changes
A simple app to gather topic changes from a Slack channel.  There's really not much to it at all... it's kind of a light wrapper around the [nlopes](https://github.com/nlopes/slack) Slack lib for Go.  You will probably need at least Go 1.11 in order to make use of the [module system](https://github.com/golang/go/wiki/Modules).

## An example
> $ ./slack-topic-changes  
> 2018-09-09 08:30:00 -0600 CST: someuser set topic to 'Oh yah, another topic'  
> 2018-09-08 08:20:00 -0600 CST: someuser set topic to 'Look, it's a new topic!'  
> 2018-09-07 08:10:00 -0600 CST: someuser set topic to 'What are topics?'  
> ...

## Steps to make this work
1. Clone the repo
2. `go build`
3. Setup environment variables

    You'll need to set two environment varibles to make this work.  SLACK_TOKEN contains an API token for the app to use.  Additional information about Slack API tokens can be found [here](https://api.slack.com/docs/token-types).  This app has been tested using user tokens (after installing an app in workspace).  Depending on your workspace setup you may need to add search:read and users:read permissions to the app as well.

    You'll also need to set SLACK_CHANNEL to the name of the channel you wish to find topic changes for.  This name is only used to construct the search query used to retrieve messages from Slack.

4. `./slack-topic-changes`