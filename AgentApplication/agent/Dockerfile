FROM openjdk:15-jdk-alpine
COPY target/agent-0.0.1-SNAPSHOT.jar agent-0.0.1-SNAPSHOT.jar
ENTRYPOINT ["java","-jar","/agent-0.0.1-SNAPSHOT.jar"]
EXPOSE 8080
