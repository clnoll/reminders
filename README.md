`reminders` uses the Temporal workflow engine to schedule and send reminders via WhatsApp.

<p align="center">
    <a href="https://swimlanes.io/u/WQXSv6BA5"><img src="https://static.swimlanes.io/27c2b46cd8322f630cdefcf7fa9ff16e.png"/></a>
<p align="center">


TODO:
- On update, fill out reminderDetails (query workflow?)
- Dismiss & update reminders via WhatsApp
- On DELETE, different message if already deleted
- Interactive reminders via child workflow
- Use continue-as-new in Workflow to keep activity count sane
- Programmatically get updated WhatsApp token
- Support YYYYMMDD HH:MM format for creating & updating reminders
- Tests for various reminder inputs
- Update parser to allow more flexibility of message content
