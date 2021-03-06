FROM openjdk:8-jre as builder
RUN apt-get update > /dev/null && apt-get install -y --no-install-recommends curl git > /dev/null
RUN mkdir minecraft
WORKDIR /minecraft
RUN curl -s -o buildtools.jar https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar
RUN java -jar buildtools.jar --rev 1.15.2

FROM golang:latest as wrapper_builder
RUN mkdir /wrapper
COPY wrapper.go /wrapper
WORKDIR /wrapper
RUN go build -o minecraft-docker-wrapper .

FROM openjdk:8-jre
RUN apt-get update && apt-get install zip
RUN mkdir /var/minecraft
COPY --from=builder /minecraft/spigot* /usr/bin
COPY --from=wrapper_builder /wrapper/minecraft-docker-wrapper /usr/bin/minecraft-docker-wrapper
WORKDIR /var/minecraft
RUN useradd -ms /bin/sh minecraft
USER minecraft
CMD ["/usr/bin/minecraft-docker-wrapper", "-f", "-Dcom.mojang.eula.agree=true", "-m", "/usr/bin/spigot-1.15.2.jar", "-j", "/usr/local/openjdk-8/bin/java"]
