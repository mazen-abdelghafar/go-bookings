#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -cache=false -production=false


# in case of using cmd flags
# -dbname= -dbuser= -dbpassword= -cache= -production=