package com.artivisi.techstack.registration.web;

import com.artivisi.techstack.registration.domain.Registration;
import com.artivisi.techstack.registration.repository.RegistrationRepository;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;
import java.time.OffsetDateTime;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import org.springframework.beans.propertyeditors.StringTrimmerEditor;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.WebDataBinder;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.InitBinder;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
public class RegistrationController {

    private final RegistrationRepository repository;

    public RegistrationController(RegistrationRepository repository) {
        this.repository = repository;
    }

    @InitBinder
    public void initBinder(WebDataBinder binder) {
        binder.registerCustomEditor(String.class, new StringTrimmerEditor(true));
    }

    @ModelAttribute("javaVersion")
    public String javaVersion() {
        return System.getProperty("java.version");
    }

    @ModelAttribute("form")
    public RegistrationForm createForm(
            @RequestParam(value = "email", required = false) String email,
            @RequestParam(value = "full_name", required = false) String fullName,
            @RequestParam(value = "phone", required = false) String phone) {
        return new RegistrationForm(email, fullName, phone);
    }

    @GetMapping("/")
    public String showForm() {
        return "form";
    }

    @PostMapping("/register")
    public String register(
            @Valid @ModelAttribute(name = "form", binding = false) RegistrationForm form,
            BindingResult result,
            HttpServletResponse response) {

        if (result.hasErrors()) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            return "form";
        }

        Registration reg = new Registration(
                UUID.randomUUID(),
                form.email().toLowerCase(),
                form.fullName(),
                form.phone(),
                OffsetDateTime.now()
        );

        try {
            repository.insert(reg);
        } catch (DuplicateKeyException e) {
            result.rejectValue("email", "duplicate", "email is already registered");
            response.setStatus(HttpStatus.CONFLICT.value());
            return "form";
        }

        return "redirect:/registrations";
    }

    @GetMapping("/registrations")
    public String list(Model model) {
        List<Registration> all = repository.findAllOrderByCreatedAtDesc();
        model.addAttribute("registrations", all);
        model.addAttribute("count", all.size());
        return "list";
    }

    @GetMapping("/health")
    @ResponseBody
    public ResponseEntity<Map<String, Object>> health() {
        try {
            repository.ping();
            return ResponseEntity.ok(Map.of("status", "ok"));
        } catch (Exception e) {
            return ResponseEntity.status(503).body(Map.of(
                    "status", "error",
                    "error", e.toString()
            ));
        }
    }
}
