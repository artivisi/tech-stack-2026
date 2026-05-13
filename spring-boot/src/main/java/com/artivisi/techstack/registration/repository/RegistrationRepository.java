package com.artivisi.techstack.registration.repository;

import com.artivisi.techstack.registration.domain.Registration;
import java.time.OffsetDateTime;
import java.util.List;
import java.util.UUID;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

@Repository
public class RegistrationRepository {

    private final JdbcTemplate jdbc;

    public RegistrationRepository(JdbcTemplate jdbc) {
        this.jdbc = jdbc;
    }

    public void insert(Registration r) {
        jdbc.update(
                "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (?, ?, ?, ?, ?)",
                r.id(),
                r.email(),
                r.fullName(),
                r.phone(),
                r.createdAt()
        );
    }

    public List<Registration> findAllOrderByCreatedAtDesc() {
        return jdbc.query(
                "SELECT id, email, full_name, phone, created_at FROM registration ORDER BY created_at DESC",
                (rs, rowNum) -> new Registration(
                        rs.getObject("id", UUID.class),
                        rs.getString("email"),
                        rs.getString("full_name"),
                        rs.getString("phone"),
                        rs.getObject("created_at", OffsetDateTime.class)
                )
        );
    }

    public void ping() {
        jdbc.queryForObject("SELECT 1", Integer.class);
    }
}
