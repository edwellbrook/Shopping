Overview
========

The goal of the project is to build an internet connected shopping list
application for a supermarket or similar store.

A customer registers their loyalty card (e.g. a nectar or tesco clubcard)
online, and enters in their weekly shopping list.

On arrival to a store they're able to pick up a shopping list device, tap their
card and the device will load in their shopping list. The device will provide an
interface for marking items in the list as collected, and to call for assistance
from a shop assistant.

If help is requested from a shop assistant, the device will broadcast it's
location to a web service that assistants can monitor. When help is requested
the device information is presented and if the customer moves about in the
store, the location will update in real time.


Project Structure
=================

The project is split up into various components and makes use of Docker to
better seperate components and reflect how they would operate in a real-world
environment.

`/app`:                An mbed application that is flashed to the nrf52
                       hardware.

                       Files:
                       - /app/src/bluetooth.h
                       - /app/src/bluetooth.cpp
                       - /app/src/display.h
                       - /app/src/display.cpp
                       - /app/src/main.cpp
                       - /app/src/nfc.h
                       - /app/src/nfc.cpp

                       Libraries:
                       - /app/C12832.lib
                       - /app/LibPN532.lib
                       - /app/mbed-os.lib

`/interface`:          A Go application that runs on a host computer interacting
                       with the IoT hardware over serial.

`/rabbitmq-web-mqtt`:  A Dockerfile and websocket plugin for running mqtt with
                       rabbitmq with docker.

`/web-customer`:       A Node.js application for the customer-facing web service
                       allowing people to register their loyalty card and
                       shopping lists.

`/web-staff`:          A static staff-facing web service for monitoring alerts
                       for customers requesting help.


Optional Categories for Marking (1)
===================================

Device control:  I have programmed a simple graphical interface for the nrf52
                 device and made use of the joystick and buttons on the mbed
                 application shield to control the interface.


