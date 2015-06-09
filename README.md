# Chicken Little
**A RESTful service to get ahold of people, quickly.**  

- Uses phone calls, SMS, and e-mail to send short messages to people registered with the service.  
- Allows for per-user configurable contact plans (e.g., "Send me an SMS.  If I don't reply within five minutes, call me on the phone.  If I don't answer, keep calling back every ten minutes until I do.").   
- Uses Twilio and Mailgun to handle the contacting.  

# Requirements
You'll need a Twilio account if you want to notify by voice and/or SMS.  You'll need a Mailgun account if you want to notify by e-mail.  Does not currently support e-mail notification by SMTP.

# API
- **[People API](https://github.com/chrissnell/chickenlittle/blob/master/docs/PEOPLE_API.md)** - used for adding and deleting people in the system.
- **[Notification Plan API](https://github.com/chrissnell/chickenlittle/blob/master/docs/NOTIFICATION_PLAN_API.md)** - used to define how people are notified (contact methods, order, and timing)
- **[Notification API](https://github.com/chrissnell/chickenlittle/blob/master/docs/NOTIFICATION_API.md)** - used to send notifications to a person using their notification plan

# Quick Start
1. You'll need [Go](http://golang.org/) installed to build the binary.

2. Fetch and build Chicken Little:
 ```
% go get github.com/chrissnell/chickenlittle
```

3. Make a directory for the config file (config.yaml) and the database (chickenlittle.db) to live:
```sudo mkdir /opt/chickenlittle```

4. Copy the binary you just built into wherever you like to keep third-party software:
```sudo cp $GOPATH/bin/chickenlittle /usr/local/bin/```

5. Copy the sample config.yaml into the directory you made in step 3:
```sudo cp $GOPATH/src/github.com/chrissnell/chickenlittle/config.yaml.sample /opt/chickenlittle/config.yaml

6. Edit the config file and fill in your Twilio and/or Mailgun API keys, endpoint URLs, etc.  For the click_url_base and callback_url_base, you can use a service like [ngrok](http://ngrok.com) for testing or you can run Chicken Little on a public network and put the base URL to your server here. 

7. Start the Chicken Little service:
```cd /opt/chickenlittle; /usr/local/bin/chickenlittle```

8. Follow the API instructions to create users and set up notification plans

# To Do
- Implement on-call scheduling to swap out notification plans depending on who is on call.
