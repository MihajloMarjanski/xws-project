package com.example.agent.model.dto;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Pattern;

public class UserCredentials {
    @NotBlank
    private String username;
    private String password;
    private Integer pin;

    public UserCredentials() {
    }

    public UserCredentials(String username, String password, Integer pin) {
        this.username = username;
        this.password = password;
        this.pin = pin;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public Integer getPin() {
        return pin;
    }

    public void setPin(Integer pin) {
        this.pin = pin;
    }
}
