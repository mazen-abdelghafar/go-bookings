#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings 


# in case of using cmd flags
# -dbname= -dbuser= -dbpassword= -cache= -production=