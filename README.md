# Spigot wrapper for minecraft

This is a docker image for spigot. It'll build spigot using the BuildTools.jar provided by spigot.

Build the spigot image by running `docker build -t mindavi/spigot .`. Then run `docker-compose up -d` to run the image. You can change the data volume in the docker-compose file.

When the service is stopped a graceful shutdown is done by sending the `stop` command to the minecraft server.

