package com.example.agent.model;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import javax.persistence.*;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Pattern;
import java.util.Collection;
import java.util.Collections;
import java.util.Date;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@MappedSuperclass
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    protected Integer id;
    @NotBlank
    protected String firstName;
    @NotBlank
    protected String lastName;
    @Column(unique=true)
    @NotBlank
    @Email
    protected String email;
    @Column(unique=true)
    @NotBlank
    protected String username;
    @NotBlank
    protected String password;
    protected boolean isActivated = true;
    protected Integer forgotten;
    protected String pin;
    protected String salt;
    protected Integer missedPasswordCounter;
    protected boolean isBlocked;
    protected Date blockedDate;
    protected Date pinCreatedDate;
}
