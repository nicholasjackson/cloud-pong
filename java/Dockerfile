FROM gradle:jdk11 AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle build --no-daemon


FROM openjdk:11-jre
LABEL author="Rosemary Wang"
COPY --from=build /home/gradle/src/build/libs/*.jar /app/cloud-pong.jar
CMD ["java", "-jar", "/app/cloud-pong.jar"]