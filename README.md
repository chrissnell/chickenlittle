# chickenlittle
**A RESTful service to get ahold of people, quickly.**  Uses phone calls, SMS, and e-mail to send short messages to people registered with the service.  Allows for per-user configurable contact plans (e.g., "Send me an SMS.  If I don't reply within five minutes, call me on the phone.  If I don't answer, keep calling back every ten minutes until I do.").   Uses Twilio and Mailgun to handle the contacting.

# Notice
This is not ready for public consumption.  Sorry.


# API
- [People API](https://github.com/chrissnell/chickenlittle/blob/master/docs/PEOPLE_API.md) - used for adding and deleting people in the system.
- [Notification Plan API](https://github.com/chrissnell/chickenlittle/blob/master/docs/NOTIFICATION_PLAN_API.md) - used to define how people are notified (contact methods, order, and timing)
