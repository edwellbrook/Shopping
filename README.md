The project is split up into various components:

 - app - this is the mbed application that is flashed to the iot hardware
 - interface - this is the go application that runs on a host computer interacting with the iot hardware over serial
 - rabbitmq-web-mqtt - this holds the dockerfile and websocket plugin for running mqtt with rabbitmq with docker
 - web-customer - this is the customer-facing web service allowing people to register their nfc card and shopping lists
 - web-staff - this is the staff-facing web service for getting alerts when a customer requests help

The project makes use of docker to better seperate components and reflect how
these components would operate in a real-world environment.
