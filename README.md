# ics2mattermost
> **Note**
>
> An ICS (calendar) parser sending todays Events and Absentees to Mattermost.

This application was created out of need for a Mattermost notification informing
us of todays appointments.
We wanted to display todays absentees, Meetingslinks and Daily times.

I use a personal access token to access our (confluence) ics calendar and parse
todays appointments. One notification gets send to our channel with all of the
aggregated information.

# Development
Setup:
```bash
pre-commit install --hook-type commit-msg
brew install go
brew install podman-compose
podman-compose up
```
Creat a new Incomig Webhook Integration and export that URL as
`MATTERMOST_URL`.

Run tests:
```bash
go tests ./...
```
