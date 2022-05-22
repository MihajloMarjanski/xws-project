package com.example.agent.model.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import javax.persistence.Column;
import javax.persistence.MappedSuperclass;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Pattern;
import java.util.Date;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class UserDto {
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
    @Pattern(regexp = "^\\S*$")
    protected String username;
    @NotBlank
    @Pattern(regexp="^(?=.*[a-zA-Z])(?=.*\\d)(?=.*[!@#$%^&*()_+\\.])[A-Za-z\\d][A-Za-z\\d!@#$%^&*()_+\\.]{8,20}$")
    protected String password;
}



