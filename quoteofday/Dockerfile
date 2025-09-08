# Simple runtime image for Java application
FROM eclipse-temurin:21-jre AS runtime

# Single layer for all setup operations
RUN groupadd -g 1001 spring && \
    useradd -u 1001 -g spring -m spring && \
    mkdir -p /app && \
    chown spring:spring /app

WORKDIR /app

# Copy the pre-built JAR file
COPY --chown=spring:spring target/quote-of-day-service-*.jar app.jar

USER spring
EXPOSE 8001

ENTRYPOINT ["java", "-jar", "app.jar"]
