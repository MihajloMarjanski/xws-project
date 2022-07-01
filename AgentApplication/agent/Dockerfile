# First stage: complete build environment
FROM maven:3.5-jdk-8-alpine AS builder

# add pom.xml and source code
ADD ./AgentApplication/agent/pom.xml pom.xml
ADD ./AgentApplication/agent/src src/

# package jar
RUN mvn clean package

# Second stage: minimal runtime environment
FROM openjdk:8-jre-alpine

WORKDIR /root/

# copy jar from the first stage
COPY --from=builder /target/agent-0.0.1-SNAPSHOT.jar agent-0.0.1-SNAPSHOT.jar

EXPOSE 8080

CMD ["java", "-jar", "agent-0.0.1-SNAPSHOT.jar"]