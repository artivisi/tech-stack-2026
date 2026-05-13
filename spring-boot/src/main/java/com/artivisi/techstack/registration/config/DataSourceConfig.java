package com.artivisi.techstack.registration.config;

import java.net.URI;
import javax.sql.DataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.jdbc.DataSourceBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class DataSourceConfig {

    @Bean
    public DataSource dataSource(@Value("${DATABASE_URL}") String databaseUrl) {
        URI uri = URI.create(databaseUrl);
        String userInfo = uri.getUserInfo();
        if (userInfo == null || !userInfo.contains(":")) {
            throw new IllegalArgumentException("DATABASE_URL must include user:password");
        }
        String[] parts = userInfo.split(":", 2);
        String host = uri.getHost();
        int port = uri.getPort();
        String dbPath = uri.getPath();
        if (host == null || port < 0 || dbPath == null || dbPath.isEmpty()) {
            throw new IllegalArgumentException("DATABASE_URL missing host/port/db path: " + databaseUrl);
        }
        String jdbcUrl = "jdbc:postgresql://" + host + ":" + port + dbPath;
        return DataSourceBuilder.create()
                .url(jdbcUrl)
                .username(parts[0])
                .password(parts[1])
                .driverClassName("org.postgresql.Driver")
                .build();
    }
}
