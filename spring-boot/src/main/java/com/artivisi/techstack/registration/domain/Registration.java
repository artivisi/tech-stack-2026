package com.artivisi.techstack.registration.domain;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.UUID;

public record Registration(
        UUID id,
        String email,
        String fullName,
        String phone,
        OffsetDateTime createdAt
) {
    private static final DateTimeFormatter UTC_PATTERN =
            DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    public String createdAtFormatted() {
        return createdAt.withOffsetSameInstant(ZoneOffset.UTC).format(UTC_PATTERN) + " UTC";
    }
}
