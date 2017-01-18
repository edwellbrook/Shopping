Overview
========

The goal of the project is to build an internet connected shopping list
application for a supermarket or similar store.

A customer registers their loyalty card (e.g. a nectar or tesco clubcard)
online, and enters in their shopping list.

On arrival to a store they're able to pick up a shopping list device, tap their
card and the device will load in their shopping list. The device will provide an
interface for marking items in the list as collected, and to call for assistance
from a shop assistant.

If help is requested from a shop assistant, the device will broadcast it's
location to a web service that assistants can monitor. When help is requested,
the device information is presented and if the customer moves about in the
store, the location will update live.



Marking Areas
=============

Completeness:     Although the plan was stripped down a little from what was
                  orignally proposed, the core ideas and technologies have still
                  been used. Within the revised plan, everything was built as I
                  wanted however if I were to continue building out the project
                  I would look into ways of allowing the mbed device to
                  communicate wirelessly with the host system instead of just
                  over serial.

Coding:           Use of Docker to seperate components to better simulate real-
                  world use. All components compile without errors. Each
                  component follows common programming patterns for the lanuage
                  used (C++, node.js, Go). Could definately use better comments
                  across the whole project.

Data Collection:  Customers are able to input their shopping list into the
                  database where it will then be transferred to the mbed device
                  when "tapping in" their NFC card.

Actuation:        When help is requested the device continually pings its
                  location providing a feed of help requests live-updating as a
                  a customer moves around the store. There is two way
                  communication between the host computer and the mbed device
                  so that the shopping list can be loaded live from the
                  database.

Visualisation:    Built two web apps (staff and customer) for composing shopping
                  list and monitoring help requests respectively. Built an
                  interactive graphical user interface on the mbed device.

Processing:       Some processing in the mbed application for filtering out
                  bluetooth beacons that aren't in a close enough range. Some
                  processing in staff web app for timing-out help requests. Some
                  processing in translating a bluetooth beacon id into a human-
                  readable location.

Video:            Uploaded to Dropbox, link at end of readme. Run-time: 7m 13s.



Project Structure & Files
=========================

The project is split up into various components and makes use of Docker to
better seperate components and reflect how they would operate in a real-world
environment.


/app
----

An mbed application that is flashed to the nrf52 hardware.

Files:

  src/bluetooth.{h,cpp}

    Functions for interacting with BLE iBeacons.

  src/display.{h,cpp}

    Functions for interacting with the mbed application shield display.

  src/main.cpp

    The main program and serial input handler.

  src/nfc.{h,cpp}

    Functions for interacting with the I2C NFC card readers.

Libraries:

  - https://developer.mbed.org/users/chris/code/C12832
  - https://developer.mbed.org/users/dotnfc/code/LibPN532
  - https://github.com/ARMmbed/mbed-os


/interface
----------

A Go application that runs on a host computer interacting with the IoT hardware
over a serial port.

Files:

  Dockerfile

    The base Dockerfile to containerise the Go application.

  src/cmd/serial-receiver/config.go

    Functions for loading application configuration.

  src/cmd/serial-receiver/database.go

    Functions for interacting with the Postgres database.

  src/cmd/serial-receiver/main.go

    Main entrypoint and processing of interactions with the mbed application
    over serial interface.

  src/cmd/serial-receiver/mqtt.go

    Functions for interacting with rabbitmq over mqtt protocol.

  src/serial/response.go

    Functions for defining and parsing the custom API for communicating over a
    serial port.

  src/serial/device.go

    Wrapper around serial port to abstract away implementation details.

Libraries:

  - https://github.com/cenkalti/backoff
  - https://github.com/lib/pq
  - https://github.com/tarm/serial
  - https://github.com/yosssi/gmq/mqtt


/rabbitmq-web-mqtt
------------------

A Dockerfile and websocket plugin for running mqtt with rabbitmq with docker.

Files:

  Dockerfile

    The base Dockerfile to containerise rabbitmq configured with MQTT and
    websockets.

Libraries:

  - https://github.com/rabbitmq/rabbitmq-web-mqtt


/web-customer
-------------

A Node.js application for the customer-facing web service allowing people to
register their loyalty card and shopping lists.

Files:

  Dockerfile

    The base Dockerfile to containerise the Node.js application.

  package.json

    Node.js project configuration file.

  app.js

    Main entrypoint to webserver

  bin/www

    Startup script

  public/stylesheets/style.css

    CSS for webpages

  routes/index.js

    Routing and request handling for webserver

  views/*.ejs

    HTML templates/responses

Libraries:

  - https://www.npmjs.com/package/bcrypt
  - https://www.npmjs.com/package/body-parser
  - https://www.npmjs.com/package/cookie-parser
  - https://www.npmjs.com/package/debug
  - https://www.npmjs.com/package/ejs
  - https://www.npmjs.com/package/express
  - https://www.npmjs.com/package/express-session
  - https://www.npmjs.com/package/morgan
  - https://www.npmjs.com/package/pg
  - https://www.npmjs.com/package/pg-pool


/web-staff
----------

A static staff-facing web service for monitoring alerts for customers requesting
help.

Files:

  Dockerfile

    The base Dockerfile to containerise the static webpage.

  index.html

    Main web page for monitoring help requests.

  javascripts/help.js

    JavaScript for connecting to rabbitmq using MQTT over websockets, handling
    help request state, and manipulating the DOM.

Libraries:

  - https://jquery.com
  - https://github.com/mqttjs/MQTT.js


Additional Files
----------------

There are additional files provided in the project to help with setup and
execution.

Files:

  db-setup.sql

    The Postgres script to set up the SQL database.

  docker-compose.yml

    A docker-compose configuration file to allow running multiple services at
    the same time. This also provides service aliases to allow referring to
    services by names rather than static IP addresses.



Optional Categories for Marking (2)
===================================

Device control:  I have programmed a simple graphical interface for the nrf52
                 device and made use of the joystick and buttons on the mbed
                 application shield to control the interface.

Robust system:   Retries with exponential backoffs have been used when
                 connecting to databases in the Go application (interface).
                 When starting all components simultatiously, sometimes the
                 databases will take longer to start up and retrying the
                 connection automatically is preferable to crashing/manually
                 waiting.



Project Video Demo
==================

A video demo of the project is available from:
https://dl.dropboxusercontent.com/u/4034363/iot-demo.mp4
