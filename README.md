# Chicken Little
**A RESTful service to get ahold of people, quickly.**  

- Uses phone calls, SMS, and e-mail to send short messages to people registered with the service.  
- Allows for per-user configurable contact plans (e.g., "Send me an SMS.  If I don't reply within five minutes, call me on the phone.  If I don't answer, keep calling back every ten minutes until I do.").   
- Uses Twilio and Mailgun (or your own SMTP server) to handle the contacting.  

# Requirements
You'll need a Twilio account if you want to notify by voice and/or SMS.  You'll need a Mailgun account or a SMTP server if you want to notify by e-mail.

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
```/usr/local/bin/chickenlittle -config PATH_TO_YOUR_CONFIG_YAML```

8. Follow the API instructions to create users and set up notification plans

# To Do
- Implement on-call rotations for teams of people
- Authentication and role-based access control (RBAC) for various API functions.
- More test coverage

# Authors
- [Christopher Snell](http://output.chrissnell.com) - Chicken Little author
- [Dominik Schulz](https://github.com/dominikschulz) - Wrote the SMTP server support and the Team API


# License
```
Copyright (C) 2015 Christopher Snell

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
