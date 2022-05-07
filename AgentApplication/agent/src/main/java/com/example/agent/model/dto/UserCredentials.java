package com.example.agent.model.dto;

public class UserCredentials {
    private String username;
    private String password;

    public UserCredentials() {
        super();
    }

    public UserCredentials(String username, String password) {
        this.setUsername(username);
        this.setPassword(password);
    }

    public String getUsername() {
        return this.username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return this.password;
    }

    public void setPassword(String password) {
        this.password = password;
    }
}
