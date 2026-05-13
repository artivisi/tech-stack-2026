package com.artivisi.techstack.registration.web;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record RegistrationForm(
        @NotBlank(message = "valid email is required")
        @Size(min = 3, max = 254, message = "email must be 3-254 characters")
        @Pattern(regexp = "^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$", message = "valid email is required")
        String email,

        @NotBlank(message = "full name is required")
        @Size(min = 2, max = 100, message = "full name must be 2-100 characters")
        @Pattern(regexp = "^[\\p{L}\\p{M}\\s.'\\-]+$", message = "full name contains invalid characters")
        String full_name,

        @NotBlank(message = "phone is required")
        @Size(min = 7, max = 20, message = "phone must be 7-20 characters")
        @Pattern(regexp = "^[+0-9 ()\\-]+$", message = "phone contains invalid characters")
        String phone
) {}
